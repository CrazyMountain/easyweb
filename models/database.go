package models

import (
	"easyweb/utils/setting"
	. "easyweb/utils/tools"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func init() {
	dbType := setting.DatabaseType
	dbUsername := setting.DatabaseUsername
	dbPassword := setting.DatabasePassword
	dbHost := setting.DatabaseHost
	dbPort := setting.DatabasePort
	dbDatabase := setting.Database

	var err error
	db, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbUsername, dbPassword, dbHost, dbPort, dbDatabase))
	CheckErr(err)

	db.DB().SetMaxOpenConns(setting.MaxOpenConn)
	db.DB().SetMaxIdleConns(setting.MaxIdleConn)
}
