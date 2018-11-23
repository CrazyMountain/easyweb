package models

import (
	"easyweb/utils/setting"
	"github.com/rs/xid"
	"time"
)

type Session struct {
	SessionID  string `gorm:"size:20;primary_key;not null;unique"`
	Username   string `gorm:"size:24"`
	ExpireTime time.Time
}

func init() {
	if !db.HasTable(&Session{}) {
		db.CreateTable(&Session{})
	}
}

func StartSession(username string) string {
	uuid := xid.New().String()
	expire := time.Now().Add(time.Duration(setting.ExpireDuration) * time.Minute)
	db.Create(&Session{uuid, username, expire})
	return uuid
}

func StopSession(id string) {
	session := Session{}
	db.Where("session_id = ?", id).First(&session).Delete(&session)
}

func IsSignIn(id string) bool {
	session := Session{}
	db.Where("session_id = ?", id).First(&session)
	if len(session.SessionID) == 0 {
		return false
	}
	if time.Now().After(session.ExpireTime) {
		db.Delete(&session)
		return false
	}
	return true
}
