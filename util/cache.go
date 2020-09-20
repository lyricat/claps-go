package util

import (
	"github.com/gin-contrib/cache/persistence"
	"github.com/spf13/viper"
	"time"
)

// 声明一个全局的rdb变量
var Rdb *persistence.RedisStore

//jwt的过期时间
const TokenExpireDuration = time.Hour* 2

// 初始化连接
func InitClient() (err error) {
	//Rdb = persistence.NewInMemoryStore(TokenExpireDuration)
	Rdb = persistence.NewRedisCache(viper.GetString("REDIS_ADDR"),
		viper.GetString("REDIS_PASSWORD"),TokenExpireDuration)

	/*
	Rdb = redis.NewClient(&redis.Options{
		Addr:     config.GetString("REDIS_ADDR"),
		Password: config.GetString("REDIS_PASSWORD"),
		DB:       config.GetInt("REDIS_DB"),  // use default DB
	})

	_, err = Rdb.Ping().Result()
	if err != nil {
		return err
	}
	 */
	return nil
}