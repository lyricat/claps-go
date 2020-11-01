package service

import (
	"claps-test/model"
	"claps-test/util"
	"context"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

/**
 * @Description: 通过bot信息创建一个mixin客户端,完成主bot初始化和对应转账功能
 * @param bot
 * @return client
 * @return err
 */
func CreateMixinClient(bot *model.Bot) (client *mixin.Client, err *util.Err) {

	s := &mixin.Keystore{
		ClientID:   bot.Id,
		SessionID:  bot.SessionId,
		PrivateKey: bot.PrivateKey,
		PinToken:   bot.PinToken,
	}

	client, err1 := mixin.NewFromKeystore(s)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrThirdParty, "创建mixinClient失败")
	}
	return
}

/**
 * @Description: 每隔5分钟异步更新数据库中的asset信息
 */
func SyncAssets() {

	ctx := context.TODO()
	for {
		assetsInfo, err := util.MixinClient.ReadAssets(ctx)
		//错误处理
		if err != nil {
			log.Error(err.Error())
			continue
		}

		for i := range assetsInfo {

			if util.CheckAsset(&assetsInfo[i].AssetID) {

				asset := &model.Asset{
					AssetId:  assetsInfo[i].AssetID,
					Symbol:   assetsInfo[i].Symbol,
					Name:     assetsInfo[i].Name,
					IconUrl:  assetsInfo[i].IconURL,
					PriceBtc: assetsInfo[i].PriceBTC,
					PriceUsd: assetsInfo[i].PriceUSD,
				}
				//第一次使用前,如果数据库没有信息,需要先创建几条记录,之后使用就每次更新即可
				err = model.ASSET.UpdateAsset(asset)

				if err != nil {
					log.Error(err.Error())
				}
			}
		}
		time.Sleep(time.Minute * 5)
	}
}

/**
 * @Description: 每隔300毫秒在数据库中获取处于未完成状态的捐赠,然后创建对应bot完成转账操作
 */
func SyncTransfer() {

	ctx := context.TODO()
	for {
		//找到状态为UNFINISHED的transfer
		transfers, err := model.TRANSFER.ListTransfersByStatus(model.UNFINISHED)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		//说明当前时间没有提现记录
		if len(*transfers) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		for _, transfer := range *transfers {
			//opponentId是转给谁
			// Transfer transfer to account
			//	asset_id, opponent_id, amount, traceID, memo
			// 把该user的钱转账到该账户返回快照

			//sender
			bot, err := model.BOT.GetBotById(transfer.BotId)
			if err != nil {
				log.Error(err.Error())
				continue
			}

			user, err1 := CreateMixinClient(bot)
			if err1 != nil {
				log.Error(err1.Errord.Error())
				continue
			}
			//traceId暂时不应该这ls
			snapshot, err := user.Transfer(ctx, &mixin.TransferInput{
				TraceID: transfer.TraceId,
				AssetID: transfer.AssetId,
				//接收方的mixin_id
				OpponentID: transfer.MixinId,
				Amount:     transfer.Amount,
				Memo:       transfer.Memo,
			}, bot.Pin)

			if err != nil {
				log.Error(err.Error())
				continue
			}

			transfer.SnapshotId = snapshot.SnapshotID
			transfer.CreatedAt = snapshot.CreatedAt
			transfer.Status = model.FINISHED

			err = model.TRANSFER.InsertOrUpdateTransfer(&transfer)
			if err != nil {
				log.Error(err.Error())
			}
		}
		time.Sleep(300 * time.Millisecond)
	}
}

/**
 * @Description: 每隔300毫秒,通过主bot异步获取mixin主网所有的转账信息,并通过判断是否有userId来判断这边转账是否针对与主bot下的子bot,
	对符合条件的的转账按照所选择的分配方式获取对应给每个member需要分配多少对应虚拟货币的金额,并加到对应member的member_wallet的total和balance字段
*/
func SyncSnapshots() {

	ctx := context.TODO()

	//获取当前时间
	since := time.Now()

	//死循环,读到上次最后一条就break
	for {
		//获取最后一次跟新记录
		property, _ := model.PROPERTY.GetPropertyByKey("last_snapshot_id")
		var lastSnapshotID string
		if property != nil {
			lastSnapshotID = property.Value
		}

		//从mixin获取当前时间之后的snapshots
		snapshots, err := util.MixinClient.ReadNetworkSnapshots(ctx, "", since, "ASC", 100)

		//错误处理
		if err != nil {
			log.Error(err.Error())
			continue
		}

		//这个时间端没有交易记录
		if len(snapshots) == 0 {
			time.Sleep(time.Second)
			continue
		}

		//遍历100记录
		for i := range snapshots {
			/*
				log.Debug(*snapshots[i])
				log.Debug("\n")
			*/
			if lastSnapshotID == snapshots[i].SnapshotID {
				continue
			}

			//筛选自己的转入
			if snapshots[i].UserID != "" && snapshots[i].Amount.Cmp(decimal.Zero) > 0 && snapshots[i].Memo != "deposit" {
				//根据机器人从数据库里找到项目
				projectTotal, err := model.PROJECTTOTAL.GetProjectTotalByBotId(snapshots[i].UserID)
				//错误处理有问题
				if err != nil {
					log.Error(err.Error())
					continue
				}

				transaction := &model.Transaction{
					Id:        snapshots[i].SnapshotID,
					ProjectId: projectTotal.Id,
					AssetId:   snapshots[i].Asset.AssetID,
					Amount:    snapshots[i].Amount,
					CreatedAt: snapshots[i].CreatedAt,
					Sender:    snapshots[i].OpponentID,
					Receiver:  snapshots[i].UserID,
				}
				//插入捐赠记录
				err = model.TRANSACTION.InsertTransaction(transaction)
				if err != nil {
					log.Error(err.Error())
					continue
				}

				//查找汇率等详细信息
				asset, err := model.ASSETTOUSD.GetPriceUsdByAssetId(snapshots[i].Asset.AssetID)
				if err != nil {
					log.Error(err.Error())
					continue
				}

				//更新Total字段
				projectTotal.Total = projectTotal.Total.Add(asset.PriceUsd.Mul(snapshots[i].Amount))
				projectTotal.Donations += 1

				err = model.PROJECTTOTAL.UpdateProjectTotal(projectTotal)
				if err != nil {
					log.Error(err.Error())
					continue
				}

				//更新项目钱包
				walletTotal, err := model.WALLETTOTAL.GetWalletTotalByBotIdAndAssetId(snapshots[i].UserID, snapshots[i].Asset.AssetID)
				if err != nil {
					log.Error(err.Error())
					continue
				}
				walletTotal.Total = walletTotal.Total.Add(snapshots[i].Amount)
				err = model.WALLET.UpdateWalletTotal(walletTotal)
				if err != nil {
					log.Error(err.Error())
					continue
				}
				//根据不同的分配算法进行配置
				bot, err := model.BOTDTO.GetBotDtoById(snapshots[i].UserID)

				switch bot.Distribution {
				case model.MericoAlgorithm:
					go distributionByMericoAlgorithm(transaction)
				case model.Commits:
					go distributionByCommits(transaction)
				case model.ChangedLines:
					go distributionByChangedLines(transaction)
				case model.IdenticalAmount:
					go distributionByIdenticalAmount(transaction)
				}
			}
		}

		lastSnapshotID = snapshots[len(snapshots)-1].SnapshotID
		property = &model.Property{
			Key:   "last_snapshot_id",
			Value: lastSnapshotID,
		}
		err = model.PROPERTY.UpdateProperty(property)
		if err != nil {
			log.Error(err.Error())
		}
		since = snapshots[len(snapshots)-1].CreatedAt
		time.Sleep(100 * time.Millisecond)
	}
}

/**
 * @Description:每隔40min,异步更新数据库中的fiat表
 */
func SyncFiat() {
	ctx := context.TODO()
	for {
		mixinFiats, err := util.MixinClient.ReadExchangeRates(ctx)
		if err != nil {
			log.Error(err.Error())
		}
		fiat := &model.Fiat{}
		for _, mixinFiat := range mixinFiats {
			fiat.Code = mixinFiat.Code
			fiat.Rate = mixinFiat.Rate
			if model.FIAT.UpdateFiat(fiat) != nil {
				log.Error(err.Error())
			}
		}

		time.Sleep(40 * time.Minute)
	}

}

/**
 * @Description: 获取认证之后的客户端
 * @param ctx
 * @param code
 * @return client
 * @return err
 */
func GetMixinAuthorizedClient(ctx *gin.Context, code string) (client *mixin.Client, err *util.Err) {
	//从配置文件中读取Id和密码
	clientId := viper.GetString("MIXIN_CLIENT_ID")
	clientSecret := viper.GetString("MIXIN_CLIENT_SECRET")

	//生成Key
	key := mixin.GenerateEd25519Key()

	//code换Token
	store, err1 := mixin.AuthorizeEd25519(ctx, clientId, clientSecret, code, "", key)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrThirdParty, "mixin Ed25519出错")
		return
	}

	//换取
	client, err2 := mixin.NewFromOauthKeystore(store)
	if err2 != nil {
		err = util.NewErr(err2, util.ErrThirdParty, "mixin store to client error")
		return
	}

	return
}

/**
 * @Description: 获取对应mixin用户信息
 * @param ctx
 * @param client
 * @return user
 * @return err
 */
func GetMixinUserInfo(ctx *gin.Context, client *mixin.Client) (user *mixin.User, err *util.Err) {

	user, err1 := client.UserMe(ctx)
	if err1 != nil {
		err = util.NewErr(err, util.ErrDataBase, "获取mixin用户的信息出错")
		return
	}
	return
}
