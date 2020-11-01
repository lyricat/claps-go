package controller

import (
	"claps-test/service"
	"claps-test/util"
	"errors"
	"github.com/gin-gonic/gin"
)

/**
 * @Description: 返回数据库中asset表中的数据
 * @param ctx
 */
func MixinAssets(ctx *gin.Context) {

	assets, err := service.ListAssetsAllByDB()
	util.HandleResponse(ctx, err, assets)
}

/**
 * @Description:  mixin oauth授权, 授权后更新数据库和缓存,由于有jwt中间件，所以一定存在cache中
 * @param ctx
 */
func MixinOauth(ctx *gin.Context) {
	type oauth struct {
		Code  string `json:"code" form:"code"`
		State string `json:"state" form:"state"`
	}

	var (
		err    *util.Err
		oauth_ oauth
	)
	resp := make(map[string]interface{})

	//获取请求参数
	if err1 := ctx.ShouldBindQuery(&oauth_); err1 != nil ||
		oauth_.Code == "" || oauth_.State == "" {
		err1 := util.NewErr(errors.New("invalid query param"), util.ErrBadRequest, "invalid query param")
		util.HandleResponse(ctx, err1, resp)
		return
	}

	var val interface{}
	var ok bool
	if val, ok = ctx.Get(util.UID); !ok {
		util.HandleResponse(ctx, util.NewErr(errors.New(""), util.ErrDataBase, "cache get uid error"), resp)
		return
	}
	uid := val.(int64)

	//用code换取令牌
	client, err := service.GetMixinAuthorizedClient(ctx, oauth_.Code)
	if err != nil {
		util.HandleResponse(ctx, err, nil)
		return
	}

	//获取mixin用户信息
	user, err2 := service.GetMixinUserInfo(ctx, client)
	if err2 != nil {
		util.HandleResponse(ctx, err2, nil)
		return
	}

	//github一定是登录,绑定mixin和github
	//更新数据库中的mixin_id字段
	err4 := service.UpdateUserMixinId(uid, user.UserID)
	if err4 != nil {
		util.HandleResponse(ctx, err4, nil)
		return
	}

	util.HandleResponse(ctx, nil, nil)
}
