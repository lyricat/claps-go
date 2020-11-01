package service

import (
	"claps-test/model"
	"claps-test/util"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

/**
 * @Description: 根据提供的不同参数,使用metric分配钱
 * @param transaction
 * @param metric
 */
func distribute(transaction *model.Transaction, metric string) {
	emailToMetric := make(map[string]decimal.Decimal)

	//Get all member by projectId
	members, err := model.USER.ListMembersByProjectId(transaction.ProjectId)
	if err != nil {
		log.Error(err.Error())
		return
	}

	//Get groupId by projectId
	groupId, err1 := GetGroupIdByProjectId(transaction.ProjectId)
	if err1 != nil {
		log.WithFields(log.Fields{
			"code": err1.Code,
			"err":  err1.Errord,
		}).Error(err1.Message)
		return
	}

	//Get devValue
	primaryEmailStrs, err := GetMetricByGroupIdAndUserEmails(groupId, metric, *members)
	if err != nil {
		log.Error("get devValue error:%v", err)
		return
	}

	if len(primaryEmailStrs) == 0 {
		log.Error("analyze not finish, use IdenticalAmount")
		distributionByIdenticalAmount(transaction)
		return
	}

	//create the map email to metric
	for _, v := range primaryEmailStrs {
		emailToMetric[v.PrimaryEmail] = v.Value
	}

	//calculate every member should get how much money
	if err1 := model.ExecuteTx(func(tx *gorm.DB) error {
		for i := range *members {
			rate := emailToMetric[(*members)[i].Email]
			amount := transaction.Amount.Mul(rate)

			//获得相应的用户钱包
			walletTotal, err := model.MEMBERWALLETDTO.GetMemberWalletByProjectIdAndUserIdAndBotIdAndAssetId(transaction.ProjectId, (*members)[i].Id, transaction.Receiver, transaction.AssetId)
			if err != nil {
				log.Error(err.Error())
				return err
			}
			walletTotal.Total = walletTotal.Total.Add(amount)
			walletTotal.Balance = walletTotal.Balance.Add(amount)

			//更新钱包
			err = model.MEMBERWALLETDTO.UpdateMemberWallet(walletTotal)
			if err != nil {
				log.Error(err.Error())
				return err
			}
		}
		return nil
	}); err1 != nil {
		err := util.NewErr(err1, util.ErrDataBase, metric+"插入提现记录事物出现问题")
		log.Error(err.Error())
		return
	}

	//自动提现的用户,在这里运行自动提现功能
	for _, v := range *members {
		err1 := WithdrawNowOrNot(&v)
		if err1 != nil {
			log.Error(err1.Error())
		}
	}
}

/**
 * @Description: 根据merico的接口获取DevValue,之后进行对对应project的members进行分配操作,并修改对应的member_wallet的balance和total字段
 * @param transaction
 */
func distributionByMericoAlgorithm(transaction *model.Transaction) {
	distribute(transaction, DEVVAL)
}

/**
 * @Description: 根据merico的接口获取对应project中members的commits值,之后进行对对应project的members进行分配操作,并修改对应的member_wallet的balance和total字段
 * @param transaction
 */
func distributionByCommits(transaction *model.Transaction) {
	distribute(transaction, COMMIT_NUM)
}

/**
 * @Description: 根据merico的接口获取对应project中members的changeLine值,之后进行对对应project的members进行分配操作,并修改对应的member_wallet的balance和total字段
 * @param transaction
 */
func distributionByChangedLines(transaction *model.Transaction) {
	distribute(transaction, CHANGE_LINES)
}

/**
 * @Description: 根据平均分配算法,之后进行对对应project的members进行分配操作,并修改对应的member_wallet的balance和total字段
 * @param transaction
 */
func distributionByIdenticalAmount(transaction *model.Transaction) {
	//通过项目ID获取所有成员
	members, err := model.USER.ListMembersByProjectId(transaction.ProjectId)
	if err != nil {
		log.Error(err.Error())
		return
	}

	//做除法,如果members等于0上面就返回
	memberNumbers := decimal.NewFromInt(int64(len(*members)))
	amount := transaction.Amount.Div(memberNumbers)

	if err1 := model.ExecuteTx(func(tx *gorm.DB) error {
		for i := range *members {
			//获得相应的用户钱包
			walletTotal, err := model.MEMBERWALLETDTO.GetMemberWalletByProjectIdAndUserIdAndBotIdAndAssetId(transaction.ProjectId, (*members)[i].Id, transaction.Receiver, transaction.AssetId)
			if err != nil {
				log.Error(err.Error())
				return err
			}
			if i == 0 {
				//因为可能会除不尽,所以这里考虑如果出现这种情况,就把除不尽的值转给第一个人
				walletTotal.Total = walletTotal.Total.Add(transaction.Amount.Sub(amount.Mul(memberNumbers.Sub(decimal.NewFromInt(1)))))
				walletTotal.Balance = walletTotal.Balance.Add(transaction.Amount.Sub(amount.Mul(memberNumbers.Sub(decimal.NewFromInt(1)))))
			} else {
				walletTotal.Total = walletTotal.Total.Add(amount)
				walletTotal.Balance = walletTotal.Balance.Add(amount)
			}
			//更新钱包
			err = model.MEMBERWALLETDTO.UpdateMemberWallet(walletTotal)
			if err != nil {
				log.Error(err.Error())
				return err
			}
		}
		return nil
	}); err1 != nil {
		err := util.NewErr(err1, util.ErrDataBase, "平均分配算法插入提现记录事物出现问题")
		log.Error(err.Error())
		return
	}

	//自动提现的用户,在这里运行自动提现功能
	for _, v := range *members {
		err1 := WithdrawNowOrNot(&v)
		if err1 != nil {
			log.Error(err1.Error())
		}
	}

}

/**
 * @Description: 立马提现或者等待用户去提现
 * @param member
 * @return err
 */
func WithdrawNowOrNot(member *model.User) (err *util.Err) {
	//判断是否有未完成的提现
	if member.WithdrawalWay == model.WithdrawByClaps {
		if member.MixinId != "" {
			err = IfUnfinishedTransfer(member.MixinId)
			if err != nil {
				return
			}
			//生成transfer记录
			err = DoTransfer(member.Id, member.MixinId)
			if err != nil {
				return
			}
		}
	}
	return
}
