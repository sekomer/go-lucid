package database

import (
	"go-lucid/config"
	"go-lucid/core/block"
	"go-lucid/core/transaction"
	"log"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func autoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&transaction.RawTransactionModel{},
		&transaction.TxInModel{},
		&transaction.TxOutModel{},
		&block.BlockModel{},
		&block.BlockHeaderModel{},
	)
}

func InitDB(path string) *gorm.DB {
	once.Do(func() {
		var err error
		db, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		config := config.MustGetFullNodeConfig()
		if config.Node.Data.AutoMigrate {
			autoMigrate(db)
		}
	})
	return db
}

func GetDB() *gorm.DB {
	return db
}

func GetTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	config := config.MustGetFullNodeConfig()
	if config.Node.Data.AutoMigrate {
		log.Println("auto migrating test db")
		autoMigrate(db)
	}
	return db
}
