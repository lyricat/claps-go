package util

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var MericoAppid string
var MericoSecret string
var GithubClinetId string
var GithubOauthCallback string
var MixinClientId string
var MixinOauthCallback string
var MySecret []byte
var Merico string

func InitConfig() {

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err.Error())
	}

	viper.SetConfigType("json")
	viper.SetConfigName(viper.GetString("MIXIN_CLIENT_CONFIG"))

	Merico = viper.GetString("MERICO_IP")
	MySecret = []byte(viper.GetString("TOKEN_SECRET"))
	//获取merci的id和secret
	MericoAppid = viper.GetString("MERICO_APPID")
	MericoSecret = viper.GetString("MERICO_SECRET")
	//获取envs
	GithubClinetId = viper.GetString("GITHUB_CLIENT_ID")
	GithubOauthCallback = viper.GetString("GITHUB_OAUTH_CALLBACK")
	MixinClientId = viper.GetString("MIXIN_CLIENT_ID")
	MixinOauthCallback = viper.GetString("MIXIN_OAUTH_CALLBACK")

	//两个配置文件合并
	err = viper.MergeInConfig()
	if err != nil {
		log.Println(err.Error())
	}

	viper.SetConfigType("yml")
	viper.SetConfigName("application.yml")
	viper.AddConfigPath("./config/")
	err = viper.MergeInConfig()
	if err != nil {
		log.Println(err.Error())
	}

}
