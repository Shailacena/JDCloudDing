package data

import (
	"apollo/server/pkg/config"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	db *gorm.DB
)

func Init(conf config.MysqlConfig) *gorm.DB {
	var err error
	dsn := conf.Uri
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database, err=%s", err))
	}
	log.Println("init db success")
	return db
}

func Instance() *gorm.DB {
	return db
}
