package dao

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	config "github.com/spf13/viper"
	"time"
	"github.com/hashicorp/go-multierror"
)

var db *gorm.DB

func InitDB() (*gorm.DB,error) {

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
		return nil,err
	}

	db.DB().SetConnMaxLifetime(time.Hour)
	db.DB().SetMaxOpenConns(1024)
	db.DB().SetMaxIdleConns(32)

	db.SingularTable(true)

	return db,nil
}

/*
获取数据库句柄
*/
func GetSqlDB() *sql.DB {
	return db.DB()
}

func ExecuteTx(fn func(*gorm.DB) error) error {
	tx := db.Begin()
	defer tx.RollbackUnlessCommitted()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit().Error
}

// Auto migrate
type MigrateHandler func(db *gorm.DB) error

var registeredMigrateHandlers []MigrateHandler

func RegisterMigrateHandler(h MigrateHandler) {
	registeredMigrateHandlers = append(registeredMigrateHandlers, h)
}

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