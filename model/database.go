package model

import (
	"database/sql"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	config "github.com/spf13/viper"
	"time"
)

var db *gorm.DB

/**
 * @Description: 初始化数据库
 * @return *gorm.DB
 * @return error
 */
func InitDB() (*gorm.DB, error) {

	if db != nil {
		return db, nil
	}
	var err error
	driverName := "postgres"
	args := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.GetString("DATABASE_HOST"),
		config.GetString("DATABASE_PORT"),
		config.GetString("DATABASE_USERNAME"),
		config.GetString("DATABASE_DATABASE"),
		config.GetString("DATABASE_PASSWORD"))

	db, err = gorm.Open(driverName, args)

	if err != nil {
		log.Panic("failed to connect database,err :" + err.Error())
		return nil, err
	}

	db.DB().SetConnMaxLifetime(time.Hour)
	db.DB().SetMaxOpenConns(1024)
	db.DB().SetMaxIdleConns(32)

	db.SingularTable(true)
	if config.GetString("gorm.logMode") == "false" {
		db.LogMode(false)
	} else {
		db.LogMode(true)
	}

	return db, nil
}

/**
 * @Description: 获取数据库句柄
 * @return *sql.DB
 */
func GetSqlDB() *sql.DB {
	return db.DB()
}

/**
 * @Description: 事物函数封装
 * @param fn
 * @return error
 */
func ExecuteTx(fn func(*gorm.DB) error) error {
	tx := db.Begin()
	defer tx.RollbackUnlessCommitted()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit().Error
}

/**
 * @Description: Auto migrate
 * @param db
 * @return error
 */
type MigrateHandler func(db *gorm.DB) error

var registeredMigrateHandlers []MigrateHandler

func RegisterMigrateHandler(h MigrateHandler) {
	registeredMigrateHandlers = append(registeredMigrateHandlers, h)
}

/**
 * @Description: model中init里注册的函数
 * @param db
 * @return error
 */
func MigrateInDB(db *gorm.DB) error {
	var result *multierror.Error
	for _, h := range registeredMigrateHandlers {
		if err := h(db); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}

func Migrate() error {
	return MigrateInDB(db)
}
