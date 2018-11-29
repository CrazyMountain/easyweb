package models

import "time"

type Model struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

type User struct {
	Model
	Username string `gorm:"size:24;not null"`
	Password string `gorm:"size:24:not null"`
}

func init() {
	// create table if not exists
	if !db.HasTable(&User{}) {
		db.CreateTable(&User{})
	}
}

func AddUser(username, password string) error {
	return db.Create(&User{Username: username, Password: password}).Error
}

func IsUserExists(username string) (bool, error) {
	if err := db.Where("username = ?", username).First(&User{}).Error; nil != err {
		return false, err
	}
	return true, nil
}

func IsPasswordCorrect(username, password string) (bool, error) {
	if err := db.Where("username = ? and password = ?", username, password).First(&User{}).Error; nil != err {
		return false, err
	}
	return true, nil
}
