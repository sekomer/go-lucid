package database

import (
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func GetDB(path string) *gorm.DB {
	once.Do(func() {
		var err error
		db, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	})
	return db
}
