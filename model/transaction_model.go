package model

import (
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"time"
)

/**
 * @Description:注册自动迁移函数
 */
func init() {
	RegisterMigrateHandler(func(db *gorm.DB) error {

		if err := db.AutoMigrate(&Transaction{}).Error; err != nil {
			return err
		}
		return nil
	})
}

type Transaction struct {
	Id        string          `json:"id,omitempty" gorm:"type:varchar(50);primary_key;not null;"`
	ProjectId int64           `json:"project_id,omitempty" gorm:"type:bigint;index:transaction_project_UNIDEX;not null"`
	AssetId   string          `json:"asset_id,omitempty" gorm:"type:varchar(50);not null"`
	Amount    decimal.Decimal `json:"amount,omitempty" gorm:"type:varchar(128);not null"`
	CreatedAt time.Time       `json:"created_at,omitempty" gorm:"type:timestamp with time zone"`
	Sender    string          `json:"sender,omitempty" gorm:"type:varchar(50);default:null"`
	Receiver  string          `json:"receiver,omitempty" gorm:"type:varchar(50);default:null"`
}

var TRANSACTION *Transaction

func (transaction *Transaction) InsertTransaction(transactionData *Transaction) (err error) {
	err = db.Create(transactionData).Error
	return
}

/**
 * @Description: 获取捐赠记录:通过projectId,已废弃
 * @receiver transaction
 * @param projectId
 * @return transactions
 * @return err
 */
func (transaction *Transaction) ListTransactionsByProjectId(projectId int64) (transactions *[]Transaction, err error) {

	transactions = &[]Transaction{}
	err = db.Debug().Where("project_id = ?", projectId).Order("created_at desc").Limit(256).Find(transactions).Error

	return
}

/**
 * @Description: 获取给项目一共获得了多少比捐赠
 * @receiver transaction
 * @param projectId
 * @return number
 * @return err
 */
func (transaction *Transaction) getTransactionsNumbersByProjectId(projectId int64) (number int, err error) {

	err = db.Debug().Table("transaction").Where("project_id = ?", projectId).Count(&number).Error
	return
}

/**
 * @Description: 通过projectId和query值,获取捐赠记录
 * @receiver transaction
 * @param projectId
 * @return transactions
 * @return err
 */
func (transaction *Transaction) ListTransactionsByProjectIdAndQuery(projectId int64, q *PaginationQ) (transactions *[]Transaction, number int, err error) {

	transactions = &[]Transaction{}
	number, err = transaction.getTransactionsNumbersByProjectId(projectId)
	if err != nil {
		return
	}

	tx := db.Debug().Table("transaction")
	if q.Limit <= 0 {
		q.Limit = 20
	}

	if q.Offset <= 0 {
		q.Offset = 0
	}
	err = tx.Where("project_id = ?",
		projectId).Order("created_at desc").Limit(q.Limit).Offset(q.Offset).Find(transactions).Error

	return
}
