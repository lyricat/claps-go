package service

import (
	"claps-test/model"
	"claps-test/util"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v32/github"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"net/http"
)

/**
 * @Description: 第一次登录没有插入,有就更新
 * @param user
 * @return err
 */
func InsertOrUpdateUser(user *model.User) (err *util.Err) {
	err1 := model.USER.InsertOrUpdateUser(user)
	if err1 != nil {
		err = util.NewErr(err, util.ErrDataBase, "数据库查询错误")
	}
	return
}

/**
 * @Description: 拼接含有code和clientID和client_secret，成一个URL用来换取Token,返回一个拼接的URL,code 表示github认证服务器返回的code
 * @param code
 * @return string
 */
func GetOauthToken(code string) string {
	str := fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		viper.GetString("GITHUB_CLIENT_ID"), viper.GetString("GITHUB_CLIENT_SECRET"), code,
	)
	//fmt.Println(str)
	return str
}

/**
 * @Description: 根据参数URL去请求，然后换取Token,返回Token指针和错误信息
 * @param url
 * @return token
 * @return err
 */
func GetToken(url string) (token *oauth2.Token, err *util.Err) {

	req, err1 := http.NewRequest(http.MethodGet, url, nil)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrInternalServer, "构建请求时发生错误")
		return
	}
	req.Header.Set("accept", "application/json")

	//发送请求并获得响应
	var httpClient = http.Client{}

	res, err2 := httpClient.Do(req)
	if err2 != nil {
		err = util.NewErr(err2, util.ErrInternalServer, "发送请求时候发生错误")
		return
	}

	//将相应体解析为token,返回
	var token1 oauth2.Token
	token = &token1

	//将返回的信息解析到Token
	if err3 := json.NewDecoder(res.Body).Decode(token); err3 != nil {
		err = util.NewErr(err3, util.ErrInternalServer, "解析Token结构体出错")
		return
	}
	log.Debug("生成的Token是", token)
	return
}

/**
 * @Description: 用获得的Token获得UserInfo,返回User指针
 * @param token
 * @return user
 * @return err
 */
func GetUserInfo(token *oauth2.Token) (user *github.User, err *util.Err) {

	log.Info("GitHub Token: ", token.AccessToken)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token.AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	user, _, err1 := client.Users.Get(ctx, "")

	if err1 != nil {
		err = util.NewErr(err1, util.ErrThirdParty, "向github请求userinfo出错")
		log.Error(err1)
		return
	}

	return
}
