package controller

import (
	"claps-test/middleware"
	"claps-test/model"
	"claps-test/service"
	"claps-test/util"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v32/github"
	log "github.com/sirupsen/logrus"
	"strconv"
)


/*
功能:用code换取Token
说明:此时没有发放token,没有cache,成功授权后发放token,设置cache
 */
func Oauth(ctx *gin.Context) {
	type oauth struct {
		Code string `json:"code" form:"code"`
	}
	var (
		err *util.Err
		oauth_ oauth
	)

	resp := make(map[string]interface{})

	if err1 := ctx.ShouldBindQuery(&oauth_);err1 != nil{
		err := util.NewErr(err1,util.ErrBadRequest, "")
		util.HandleResponse(ctx, err, resp)
		return
	}
	log.Debug("code = ",oauth_.Code)

	var oauthTokenUrl = service.GetOauthToken(oauth_.Code)
	//处理请求的URL,获得Token指针
	token2,err := service.GetToken(oauthTokenUrl)
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

	log.Debug("user = ",*user)

	//生成token
	jwt_token,err1 := middleware.GenToken(strconv.FormatInt(*user.ID,10))
	if err1 != nil{
		util.HandleResponse(ctx,util.NewErr(err1,util.ErrInternalServer,"gen token error"),nil)
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
	for _,v := range emails{
		if *v.Primary{
			u.Email = *v.Email
			break
		}
	}


	err = service.InsertOrUpdateUser(&u)
	if err != nil {
		util.HandleResponse(ctx, err, resp)
		return
	}

	//redis存储user信息
	mcache := &util.MCache{}
	emailForCache := []github.UserEmail{}
	for _,val:= range emails{
		emailForCache = append(emailForCache, *val)
	}
	mcache.Github = *user
	mcache.GithubEmails = emailForCache
	mcache.GithubAuth = true

	err1 = util.Rdb.Set(strconv.FormatInt(*user.ID,10),mcache,-1)
	if err1 != nil{
		util.HandleResponse(ctx,util.NewErr(err1,util.ErrDataBase,"set cache error"),nil)
		return
	}

	resp["token"] = jwt_token
	util.HandleResponse(ctx,nil,resp)
	//重定向到http://localhost:3000/profile
	//newpath := "http://localhost:3000" + oauth_.Path
	//log.Debug("重定向", newpath)
	//ctx.Redirect(http.StatusMovedPermanently, newpath)

}
