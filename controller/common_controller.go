package controller

import (
	"claps-test/middleware"
	"claps-test/model"
	"claps-test/service"
	"claps-test/util"
	"errors"
	"github.com/gin-gonic/gin"
)

/**
 * @Description: 用code换取Token,此时没有发放token,成功授权后发放token,token中记录github的userId
 * @param ctx
 */
func Oauth(ctx *gin.Context) {
	type oauth struct {
		Code string `json:"code" form:"code"`
	}
	var (
		err    *util.Err
		oauth_ oauth
	)

	resp := make(map[string]interface{})

	if err1 := ctx.ShouldBindQuery(&oauth_); err1 != nil {
		err := util.NewErr(err1, util.ErrBadRequest, "")
		util.HandleResponse(ctx, err, resp)
		return
	}

	var oauthTokenUrl = service.GetOauthToken(oauth_.Code)
	//处理请求的URL,获得Token指针
	token2, err := service.GetToken(oauthTokenUrl)
	if err != nil {
		util.HandleResponse(ctx, err, resp)
		return
	}

	// 通过token，获取用户信息
	user, err := service.GetUserInfo(token2)
	if err != nil {
		util.HandleResponse(ctx, err, resp)
		return
	}

	//通过token,获取Email信息
	emails, err := service.ListEmailsByToken(token2.AccessToken)
	//如果因为超时出错,重新请求
	if err != nil {
		util.HandleResponse(ctx, err, resp)
		return
	}

	//生成token
	//randomUid := util.RandUp(32)
	randomUid := *user.ID
	jwtToken, err1 := middleware.GenToken(randomUid)
	if err1 != nil {
		util.HandleResponse(ctx, util.NewErr(err1, util.ErrInternalServer, "gen token error"), nil)
		return
	}


	//向数据库中存储user信息
	u := model.User{}
	u.Id = *user.ID
	u.Name = *user.Login
	if user.AvatarURL != nil {
		u.AvatarUrl = *user.AvatarURL
	}
	if user.Name != nil {
		u.DisplayName = *user.Name
	}
	for _, v := range emails {
		//主email,参与分钱
		if *v.Primary {
			u.Email = *v.Email
			break
		}
	}

	//每次授权后都更新数据库中的信息
	err = service.InsertOrUpdateUser(&u)
	if err != nil {
		util.HandleResponse(ctx, err, resp)
		return
	}

	//token 的uid是github的userId
	resp["token"] = jwtToken
	util.HandleResponse(ctx, nil, resp)
}

/**
 * @Description: 返回环境信息,此时用户没有登录github没有
 * @param ctx
 */
func Environments(ctx *gin.Context) {
	resp := make(map[string]interface{})

	resp["GITHUB_CLIENT_ID"] = util.GithubClinetId
	resp["GITHUB_OAUTH_CALLBACK"] = util.GithubOauthCallback
	resp["MIXIN_CLIENT_ID"] = util.MixinClientId
	resp["MIXIN_OAUTH_CALLBACK"] = util.MixinOauthCallback

	util.HandleResponse(ctx, nil, resp)
}

/**
 * @Description: 认证用户信息,判断github和mixin是否登录绑定,之前有JWTAuthMiddleWare,有jwt说明一定github授权,ctx里设置uid
 * @param ctx
 */
func AuthInfo(ctx *gin.Context) {
	resp := make(map[string]interface{})

	var val interface{}
	var ok bool
	if val, ok = ctx.Get(util.UID); !ok {
		util.HandleResponse(ctx, util.NewErr(errors.New(""), util.ErrDataBase, "ctx get uid error"), resp)
		return
	}
	uid := val.(int64)

	var mixinAuth bool
	mixinId,err := service.GetMixinIdByUserId(uid)
	if err != nil{
		util.HandleResponse(ctx,err,resp)
		return
	}
	//没有绑定mixin
	if mixinId != ""{
		mixinAuth = true
	}else {
		mixinAuth = false
	}

	user,err := service.GetUserById(uid)
	if err != nil{
		util.HandleResponse(ctx,err,resp)
		return
	}

	resp["user"] = *user
	resp["randomUid"] = uid
	resp["mixinAuth"] = mixinAuth

	util.HandleResponse(ctx, nil, resp)
	return
}

