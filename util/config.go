package util

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

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
