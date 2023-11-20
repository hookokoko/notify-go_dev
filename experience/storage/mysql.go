package storage

import (
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var once *sync.Once

func MysqlDB() *gorm.DB {
	var db *gorm.DB
	once.Do(func() {
		db = newMysqlDb()
	})
	return db
}

func newMysqlDb() *gorm.DB {
	dsn := "root:@tcp(127.0.0.1:3306)/message_go?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
