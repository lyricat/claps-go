package controller

import (
	"claps-test/service"
	"claps-test/util"
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func MixinAssets(ctx *gin.Context) {

	assets, err := service.ListAssetsAllByDB()
	util.HandleResponse(ctx, err, assets)
}

/*
功能:mixin oauth授权
说明:授权后更新数据库和缓存
 */
func MixinOauth(ctx *gin.Context) {
	type oauth struct {
		Code string `json:"code" form:"code"`
		State string `json:"state" form:"state"`
	}

	var (
		err *util.Err
		oauth_ oauth
		//randomUid = ""
	)
	resp := make(map[string]interface{})

	//获取请求参数
	if err1 := ctx.ShouldBindQuery(&oauth_);err1 != nil ||
		oauth_.Code =="" || oauth_.State == "" {
		err1 := util.NewErr(errors.New("invalid query param"), util.ErrBadRequest, "")
		util.HandleResponse(ctx, err1, resp)
		return
	}
	log.Debug("code = ",oauth_.Code)
	log.Debug("state = ",oauth_.State)


	mcache := &util.MCache{}
	err1 := util.Rdb.Get(oauth_.State,mcache)
	//验证state
	if err1 != nil{
		err = util.NewErr(err1,util.ErrBadRequest, "invalid oauth state")
		util.HandleResponse(ctx, err, resp)
		return
	}

	//用state换取令牌
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

	//更新cache
	mcache.MixinAuth = true
	mcache.MixinId = user.UserID
	err1 = util.Rdb.Replace(oauth_.State,mcache,-1)
	if err1 != nil{
		err = util.NewErr(errors.New("cache error"), util.ErrDataBase, "")
		util.HandleResponse(ctx, err, resp)
		return
	}

	log.Debug("update mixin_id by user_id")
	//github一定是登录,绑定mixin和github
	//更新数据库中的mixin_id字段
	err4 := service.UpdateUserMixinId(*mcache.Github.ID, user.UserID)
	if err4 != nil {
		util.HandleResponse(ctx, err4, nil)
		return
	}

	//重定向
	ctx.Redirect(http.StatusMovedPermanently, "http://localhost:3000/assets")
}
