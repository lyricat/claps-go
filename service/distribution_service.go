package service

import (
	"claps-test/dao"
	"claps-test/model"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func distributionByPersperAlgorithm(transaction *model.Transaction) {

}

func distributionByCommits(transaction *model.Transaction) {

}

func distributionByChangedLines(transaction *model.Transaction) {

}

//平均分配算法
func distributionByIdenticalAmount(transaction *model.Transaction) {
	//通过项目ID获取所有成员
	members, err := dao.ListMembersByProjectId(transaction.ProjectId)
	if err != nil {
		log.Error(err.Error())
		return
	}

	//做除法,如果members等于0上面就返回?
	memberNumbers := decimal.NewFromInt(int64(len(*members)))
	amount := transaction.Amount.Div(memberNumbers)
	for i := range *members {
		//获得相应的用户钱包
		walletTotal, err := dao.GetMemberWalletByProjectIdAndUserIdAndBotIdAndAssetId(transaction.ProjectId, (*members)[i].Id, transaction.Receiver, transaction.AssetId)
		if err != nil {
			log.Error(err.Error())
			return
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
		err = dao.UpdateMemberWallet(walletTotal)
		if err != nil {
			log.Error(err.Error())
			return
		}
	}
}
