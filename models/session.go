package models

import (
	"easyweb/utils/setting"
	"github.com/rs/xid"
	"time"
)

type Session struct {
	UUID       string `gorm:"size:20"`
	Username   string `gorm:"size:24"`
	ExpireTime time.Time
}

func init() {
	if !db.HasTable(&Session{}) {
		db.CreateTable(&Session{})
	}
}

func StartSession(username string) {
	uuid := xid.New().String()
	expire := time.Now().Add(time.Duration(setting.ExpireDuration) * time.Minute)
	db.Create(&Session{uuid, username, expire})
}
