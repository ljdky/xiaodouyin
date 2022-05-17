package db

import (
	"xiaodouyin/pkg/constants"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	DB, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               constants.DSN,
		DefaultStringSize: 512,
	}))

	if err != nil {
		panic(err.Error())
	}

	err = DB.AutoMigrate(&Video{}, &Favorite{}, &Comment{})

	if err != nil {
		panic(err.Error())
	}
}
