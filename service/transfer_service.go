package service

import (
	"claps-test/dao"
	"claps-test/model"
	"claps-test/util"
	"errors"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func InsertTransfer(botId, assetID, memo string, amount decimal.Decimal, mixinId string) (err error) {
	transfer := &model.Transfer{
		BotId:   botId,
		MixinId: mixinId,
		TraceId: mixinId + assetID,
		//TransactionId: TransactionID,
		AssetId: assetID,
		Amount:  amount,
		Memo:    memo,
		Status:  model.UNFINISHED,
	}

	err1 := dao.InsertTransfer(transfer)

	if err1 != nil {
		err = util.NewErr(err, util.ErrDataBase, "根据分配算法分配之后算出应为每位member分配多少钱,提现记录首次写入数据库失败导致提现失败")
	}
	return
}

//跟新数据库中的transger表中的status为1
func UpdateTransferStatusByAssetIdAndUserId(mixinId string, assetId string, status string) (err *util.Err) {
	err1 := dao.UpdateTransferStatusByUserIdAndAssetId(mixinId, assetId, status)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "更新数据库transfer状态出错")
	}
	return
}

//判断某种币是否有未完成的提现操作,err非nil标有有未完成,err=nil表示没有未完成
func IfUnfinishedTransfer(mixinId string) (err *util.Err) {
	count, err1 := dao.CountUnfinishedTransfer(mixinId)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "数据库查询Unfinished出错")
		return
	}
	if count != 0 {
		err = util.NewErr(errors.New("有未完成的提现操作"), util.ErrForbidden, "有未完成的提现操作")
		return
	}

	log.Debug(count)
	return
}

//生成trasfer记录
func DoTransfer(userId int64, mixinId string) (err *util.Err) {

	memberWallets, err1 := dao.GetMemeberWalletByUserId(userId)
	if err1 != nil {
		err = util.NewErr(err, util.ErrDataBase, "获取用户钱包失败导致提现失败")
		return
	}

	for _,value := range *memberWallets {
		if !value.Balance.Equal(decimal.Zero) {
			err1 = InsertTransfer(value.BotId, value.AssetId, "恭喜您获得一笔捐赠", value.Balance, mixinId)
			if err1 != nil {
				err = util.NewErr(err1, util.ErrDataBase, "插入提现记录失败")
				return
			}

			//清零
			value.Balance = decimal.Zero
			//更新member_wallet
			err1 = dao.UpdateMemberWallet(&value)
			if err1 != nil {
				err = util.NewErr(err1, util.ErrDataBase, "更新用户钱包可提现值导致提现失败")
				return
			}
		}
	}

	return
}
