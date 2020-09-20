package service

import (
	"claps-test/dao"
	"claps-test/model"
	"claps-test/util"
	"context"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

//获取所有币的信息
func GetAssetByMixinClient(botId string, assetId string) (asset *mixin.Asset, err *util.Err) {
	bot, err1 := dao.GetBotById(botId)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "通过相应botid获取bot信息信息错误")
		return
	}
	mixinClient, err := CreateMixinClient(bot)
	if err != nil {
		return
	}
	asset, err1 = mixinClient.ReadAsset(context.Background(), assetId)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrThirdParty, "通过mixin获取asset信息错误")
	}
	return

}

func CreateMixinClient(bot *model.Bot) (client *mixin.Client, err *util.Err) {

	s := &mixin.Keystore{
		ClientID:   bot.Id,
		SessionID:  bot.SessionId,
		PrivateKey: bot.PrivateKey,
		PinToken:   bot.PinToken,
	}

	client, err1 := mixin.NewFromKeystore(s)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrThirdParty, "创建mixinclient失败")
	}
	return
}

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
				//err = dao.InsertAsset(asset)
				err = dao.UpdateAsset(asset)

				if err != nil {
					log.Error(err.Error())
				}
			}
		}
		time.Sleep(time.Minute * 5)
	}
}

func SyncTransfer() {

	ctx := context.TODO()
	for {
		//找到状态为UNFINISHED的trasfer
		transfers, err := dao.ListTransfersByStatus(model.UNFINISHED)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		//说明当前时间没有提现记录
		if len(*transfers) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		for i := range *transfers {
			//opponentid是转给谁
			// Transfer transfer to account
			//	asset_id, opponent_id, amount, traceID, memo
			// 把该user的钱转账到该账户返回快照

			//sender
			bot, err := dao.GetBotById((*transfers)[i].BotId)
			if err != nil {
				log.Error(err.Error())
				continue
			}

			user, err1 := CreateMixinClient(bot)
			if err1 != nil {
				log.Error(err1.Errord.Error())
				continue
			}
			//traceid暂时不应该这ls
			snapshot, err := user.Transfer(ctx, &mixin.TransferInput{
				TraceID: uuid.Must(uuid.NewV4()).String(),
				AssetID: (*transfers)[i].AssetId,
				//接收方的mixin_id
				OpponentID: (*transfers)[i].MixinId,
				Amount:     (*transfers)[i].Amount,
				Memo:       (*transfers)[i].Memo,
			}, bot.Pin)

			if err != nil {
				log.Error(err.Error())
				continue
			}

			transfer := &map[string]interface{}{
				"status":      model.FINISHED,
				"trace_id":    snapshot.TraceID,
				"snapshot_id": snapshot.SnapshotID,
				"created_at":  snapshot.CreatedAt,
			}
			//log.Error(transfer["status"])
			//更新trace_id为随机数,主键改变了，不能save
			err = dao.UpdateTransferTraceId(transfer, (*transfers)[i].TraceId)

			if err != nil {
				log.Error(err.Error())
			}
		}
		time.Sleep(300 * time.Millisecond)
	}
}

func SyncSnapshots() {

	ctx := context.TODO()

	//获取当前时间
	since := time.Now()

	//死循环,读到上次最后一条就break
	for {
		//获取最后一次跟新记录
		property, _ := dao.GetPropertyByKey("last_snapshot_id")
		var lastSnapshotID string
		if property != nil{
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
			if snapshots[i].UserID != "" && snapshots[i].Amount.Cmp(decimal.Zero) > 0 {
				//根据机器人从数据库里找到项目
				projectTotal, err := dao.GetProjectTotalByBotId(snapshots[i].UserID)
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
				err = dao.InsertTransaction(transaction)
				if err != nil {
					log.Error(err.Error())
					continue
				}

				//查找汇率等详细信息
				asset, err := dao.GetPriceUsdByAssetId(snapshots[i].Asset.AssetID)
				if err != nil {
					log.Error(err.Error())
					continue
				}

				//更新Total字段
				projectTotal.Total = projectTotal.Total.Add(asset.PriceUsd.Mul(snapshots[i].Amount))
				projectTotal.Donations += 1

				err = dao.UpdateProjectTotal(projectTotal)
				if err != nil {
					log.Error(err.Error())
					continue
				}

				//更新项目钱包
				walletTotal, err := dao.GetWalletTotalByBotIdAndAssetId(snapshots[i].UserID, snapshots[i].Asset.AssetID)
				if err != nil {
					log.Error(err.Error())
					continue
				}
				walletTotal.Total = walletTotal.Total.Add(snapshots[i].Amount)
				err = dao.UpdateWalletTotal(walletTotal)
				if err != nil {
					log.Error(err.Error())
					continue
				}
				//根据不同的分配算法进行配置
				bot, err := dao.GetBotDtoById(snapshots[i].UserID)

				switch bot.Distribution {
				case model.PersperAlgorithm:
					distributionByPersperAlgorithm(transaction)
				case model.Commits:
					distributionByCommits(transaction)
				case model.ChangedLines:
					distributionByChangedLines(transaction)
				case model.IdenticalAmount:
					distributionByIdenticalAmount(transaction)
				}
			}
		}

		lastSnapshotID = snapshots[len(snapshots)-1].SnapshotID
		property = &model.Property{
			Key:   "last_snapshot_id",
			Value: lastSnapshotID,
		}
		err = dao.UpdateProperty(property)
		if err != nil {
			log.Error(err.Error())
		}
		since = snapshots[len(snapshots)-1].CreatedAt
		time.Sleep(100 * time.Millisecond)
	}
}

//获取认证之后的客户端
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

	/*
		//将client存入session
		session := sessions.Default(ctx)
		session.Set("mixinClient",client)
		err3 := session.Save()
		if err3 != nil{
			err = util.NewErr(err3,util.ErrInternalServer,"设置mixin client session 出错")
			return
		}
	*/

	return
}

func GetMixinUserInfo(ctx *gin.Context, client *mixin.Client) (user *mixin.User, err *util.Err) {

	user, err1 := client.UserMe(ctx)
	if err1 != nil {
		err = util.NewErr(err, util.ErrDataBase, "获取mixin用户的信息出错")
		return
	}
	return
}
