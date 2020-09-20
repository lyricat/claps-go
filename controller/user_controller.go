package controller

import (
	"claps-test/service"
	"claps-test/util"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

/*
功能:返回用户的邮箱数组和项目数组
说明:执行该函数,一定登录了github
 */
func UserProfile(ctx *gin.Context) {
	var err *util.Err
	resp := make(map[string]interface{})

	//获取github信息
	mcache := util.MCache{}
	uid := ctx.GetString(util.UID)
	err1 := util.Rdb.Get(uid,&mcache)
	if err1 != nil{
		util.HandleResponse(ctx,util.NewErr(err1,util.ErrDataBase,"Redis get error."),nil)
		return
	}
	log.Debug("mcache = ",mcache)


	//根据userId获取所有project信息,Total和Patrons字段添加
	projects, err := service.ListProjectsByUserId(*mcache.Github.ID)
	if err != nil {
		util.HandleResponse(ctx, err, resp)
		return
	}

	resp["emails"] = mcache.GithubEmails
	resp["projects"] = projects
	util.HandleResponse(ctx, err, resp)
}

/*
功能:获取用户钱包所有币种的余额
说明:此时已经登录github,不需要绑定mixin
 */
func UserAssets(ctx *gin.Context) {

	var err *util.Err
	resp := make(map[string]interface{})

	var val interface{}
	var ok bool
	if val,ok = ctx.Get(util.UID);!ok{
		util.HandleResponse(ctx,util.NewErr(errors.New(""),util.ErrDataBase,"cache get uid error"),resp)
		return
	}
	uid := val.(string)

	//从redis获取cache
	mcache := &util.MCache{}
	err1 := util.Rdb.Get(uid,mcache)
	if err1 != nil{
		util.HandleResponse(ctx,util.NewErr(err1,util.ErrDataBase,"cache get error"),resp)
		return
	}

	//获得所有币的信息
	assets, err := service.ListAssetsAllByDB()
	if err != nil {
		util.HandleResponse(ctx, err, resp)
		return
	}

	//查询用户钱包,获得相应的余额,添加到币信息的后面
	err2, dto := service.GetUserBalanceByAllAssets(*mcache.Github.ID, assets)
	if err2 != nil {
		util.HandleResponse(ctx, err, resp)
		return
	}

	resp["assets"] = dto
	util.HandleResponse(ctx, err, resp)
}

//从transaction中读取关于自己项目的所有捐赠
func UserTransactions(ctx *gin.Context) {
	resp := make(map[string]interface{})

	assetId := ctx.Query("assetId")
	if assetId == "" {
		err := util.NewErr(nil, util.ErrBadRequest, "没有币种参数")
		util.HandleResponse(ctx, err, resp)
		return
	}
	log.Debug("assetId = ", assetId)

	//从transfer表中获取该用户的所有捐赠记录
}

/*
功能:获取某币种的交易记录,从transfer表里面读取数据
说明:用户一定登录了github和mixin,中间件保证
 */
func UserTransfer(ctx *gin.Context) {
	resp := make(map[string]interface{})

	var val interface{}
	var ok bool
	if val,ok = ctx.Get(util.UID);!ok{
		util.HandleResponse(ctx,util.NewErr(errors.New(""),util.ErrDataBase,"ctx get uid error"),resp)
		return
	}
	uid := val.(string)

	mcache := &util.MCache{}
	err1 := util.Rdb.Get(uid,mcache)
	if err1 != nil{
		util.HandleResponse(ctx,util.NewErr(err1,util.ErrDataBase,"cache get error"),resp)
		return
	}

	//用户如果提现过一定是绑定了mixin,没有mixin则是没有提现记录
	mixinId, err := service.GetMixinIdByUserId(*mcache.Github.ID)
	if err != nil {
		util.HandleResponse(ctx, err, nil)
		return
	}

	//可能封装成中间件,判断mixin是否登录
	if mixinId == "" {
		util.HandleResponse(ctx, util.NewErr(nil, util.ErrUnauthorized, "没有绑定mixin没有提现记录"), nil)
		return
	}

	//从transfer表中获取该用户的所有捐赠记录
	transfers, err := service.GetTransferByMininId(mixinId)
	resp["transfers"] = transfers
	util.HandleResponse(ctx, err, resp)
}

//获取某用户的所有的受捐赠记录的汇总
/*
功能:请求获得某个用户的捐赠信息的汇总,包括总金额和捐赠人数
说明:不需要绑定mixin
 */
func UserDonation(ctx *gin.Context) {
	resp := make(map[string]interface{})

	var val interface{}
	var ok bool
	if val,ok = ctx.Get(util.UID);!ok{
		util.HandleResponse(ctx,util.NewErr(errors.New(""),util.ErrDataBase,"ctx get uid error"),resp)
		return
	}
	uid := val.(string)

	mcache := &util.MCache{}
	err1 := util.Rdb.Get(uid,mcache)
	if err1 != nil{
		util.HandleResponse(ctx,util.NewErr(err1,util.ErrDataBase,"cache get error"),resp)
		return
	}

	//读取所有的member_wallet表然后汇总
	//获得所有币的信息
	assets, err := service.ListAssetsAllByDB()
	if err != nil {
		util.HandleResponse(ctx, err, resp)
		return
	}
	log.Debug(*assets)

	//查询用户钱包,获得相应的余额,添加到币信息的后面
	err2, dto := service.GetUserBalanceByAllAssets(*mcache.Github.ID, assets)
	if err2 != nil {
		util.HandleResponse(ctx, err, resp)
		return
	}

	//便利dto然后求和
	var sum decimal.Decimal

	for i := range *dto {
		sum = sum.Add((*dto)[i].Total)
	}

	//从project里面寻找Donations然后求和
	donations, err3 := service.SumProjectDonationsByUserId(*mcache.Github.ID)
	if err3 != nil {
		util.HandleResponse(ctx, err3, resp)
		return
	}

	log.Debug("donations = ",donations)

	resp["total"] = sum
	resp["donations"] = donations

	util.HandleResponse(ctx, nil, resp)
}

/*
功能:用户提现某种货币,把表中的status由0变为1
说明:已经中间件验证绑定了mixin
 */
func UserWithdraw(ctx *gin.Context) {
	resp := make(map[string]interface{})

	var val interface{}
	var ok bool
	if val,ok = ctx.Get(util.UID);!ok{
		util.HandleResponse(ctx,util.NewErr(errors.New(""),util.ErrDataBase,"ctx get uid error"),resp)
		return
	}
	uid := val.(string)

	mcache := &util.MCache{}
	err1 := util.Rdb.Get(uid,mcache)
	if err1 != nil{
		util.HandleResponse(ctx,util.NewErr(err1,util.ErrDataBase,"cache get error"),resp)
		return
	}

	//已经绑定mixin,直接从缓存中取
	mixinId := mcache.MixinId

	//判断是否有未完成的提现
	err3 := service.IfUnfinishedTransfer(mixinId)
	if err3 != nil {
		util.HandleResponse(ctx, err3, nil)
		return
	}

	//生成trasfer记录
	err2 := service.DoTransfer(*mcache.Github.ID, mixinId)
	if err2 != nil {
		util.HandleResponse(ctx, err2, nil)
		return
	}
	util.HandleResponse(ctx, nil, nil)

	//等协程完成转账
}
