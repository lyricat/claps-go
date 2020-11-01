package service

import (
	"claps-test/model"
	"claps-test/util"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

/**
 * @Description: 为避免提现多次转账,先首次生成提现记录到数据库,提现记录状态为unfinished,如果某一个用户有unfinished的记录则不允许提现,
每隔300毫秒,数据库中会异步获取为完成的提现记录,真正调用转账函数,如果转账成功,则修改转账记录状态为finished,否则不修改状态值,在此期间用户不能二次提现
 * @param botId
 * @param assetID
 * @param memo
 * @param amount
 * @param mixinId
 * @return err
*/
func InsertTransfer(botId, assetID, memo string, amount decimal.Decimal, mixinId string) (err error) {
	transfer := &model.Transfer{
		BotId:   botId,
		MixinId: mixinId,
		TraceId: uuid.Must(uuid.NewV4()).String(),
		AssetId: assetID,
		Amount:  amount,
		Memo:    memo,
		Status:  model.UNFINISHED,
	}

	err = model.TRANSFER.InsertOrUpdateTransfer(transfer)
	if err != nil {
		log.Error("model.InsertOrUpdateTransfer 错误", err)
	}

	return
}

/**
 * @Description: 判断某种币是否有未完成的提现操作,err非nil标有有未完成,err=nil表示没有未完成
 * @param mixinId
 * @return err
 */
func IfUnfinishedTransfer(mixinId string) (err *util.Err) {
	count, err1 := model.TRANSFER.CountUnfinishedTransfer(mixinId)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "数据库查询Unfinished出错")
		return
	}
	if count != 0 {
		err = util.NewErr(errors.New("该提现用户有未完成的提现操作"), util.ErrForbidden, "该提现用户有未完成的提现操作")
		return
	}

	return
}

/**
 * @Description: 生成transfer记录
 * @param userId
 * @param mixinId
 * @return err
 */
func DoTransfer(userId int64, mixinId string) (err *util.Err) {

	memberWallets, err1 := model.MEMBERWALLETDTO.GetMemberWalletByUserId(userId)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "获取用户钱包失败导致提现失败")
		return
	}

	if err1 := model.ExecuteTx(func(tx *gorm.DB) error {
		for _, value := range *memberWallets {
			if !value.Balance.Equal(decimal.Zero) {
				err2 := InsertTransfer(value.BotId, value.AssetId, "Congratulations on your donation!", value.Balance, mixinId)
				if err2 != nil {
					err = util.NewErr(err2, util.ErrDataBase, "插入提现记录失败")
					return err
				}
			}
		}
		return nil
	}); err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "插入提现记录事物出现问题")
		return err
	}

	//更新member_wallet
	err1 = model.MEMBERWALLET.UpdateMemberWalletBalanceToZeroByUserId(userId)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "更新用户钱包可提现值导致提现失败")
		return
	}
	return
}
