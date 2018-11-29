package models

import (
	"easyweb/utils/setting"
	"github.com/rs/xid"
	"time"
)

type SessionModel struct {
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

type Session struct {
	SessionModel
	SessionID string    `gorm:"size:20;primary_key;not null;unique"`
	Username  string    `gorm:"size:24"`
	ExpiredAt time.Time `gorm:"not null"`
}

func init() {
	// create table if not exists
	if !db.HasTable(&Session{}) {
		db.CreateTable(&Session{})
	}
}

func StartSession(username string) (string, error) {
	uuid := xid.New().String()
	expiredAt := time.Now().Add(time.Duration(setting.ExpireDuration) * time.Minute)
	err := db.Create(&Session{SessionID: uuid, Username: username, ExpiredAt: expiredAt}).Error
	return uuid, err
}

func EndSession(id string) error {
	return db.Where("session_id = ?", id).First(&Session{}).Delete(&Session{}).Error
}

func IsSignIn(id string) (bool, error) {
	session := Session{}
	err := db.Where("session_id = ?", id).First(&session).Error
	if nil != err {
		return false, err
	}
	if time.Now().After(session.ExpiredAt) {
		return false, EndSession(id)
	}
	return true, nil
}
