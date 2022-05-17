package db

import (
	"xiaodouyin/pkg/constants"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:               constants.DSN,
		DefaultStringSize: 512, // 默认字符串最长长度512
	}), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	err = DB.AutoMigrate(&User{}, &Follow{})

	if err != nil {
		panic(err.Error())
	}
}
