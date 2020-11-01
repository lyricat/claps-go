package main

import (
	"claps-test/model"
	"claps-test/router"
	"claps-test/service"
	"claps-test/util"
	"flag"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

/**
 * @Description:初始化配置文件,Mixin,log,DB和cache
 */
func initAllConfig() {
	util.InitConfig()
	util.InitMixin()
	util.InitLog()
}

func main() {

	cmd := flag.String("cmd", "", "process identity")
	flag.Parse()
	initAllConfig()
	db, _ := model.InitDB()
	if db != nil {
		defer db.Close()
	}
	switch *cmd {
	case "migrate", "setdb":
		if multierror := model.Migrate(); multierror != nil {
			log.Error(multierror)
		}
	default:
		//定期更新数据库snapshot信息
		go service.SyncSnapshots()
		//定期更新数据库asset信息
		go service.SyncAssets()
		//定期进行提现操作,并更改数据库
		go service.SyncTransfer()
		//定期获取汇率
		go service.SyncFiat()

		r := gin.New()
		r = router.CollectRoute(r)
		serverPort := viper.GetString("server.port")
		if serverPort != "" {
			panic(r.Run(":" + serverPort))
		} else {
			panic(r.Run(":3001"))
		}

	}
}
