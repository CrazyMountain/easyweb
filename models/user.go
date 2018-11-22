package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:24"`
	Password string `gorm:"size:24"`
}

func init() {
	if !db.HasTable(&User{}) {
		db.CreateTable(&User{})
	}
}

func AddUser(username, password string) {
	model := gorm.Model{CreatedAt: time.Now()}
	user := User{model, username, password}
	db.Create(&user)
}

func IsUserExists(username string) bool {
	var user User
	db.Where("username = ?", username).First(&user)
	fmt.Println(user)
	return user.ID != 0
}

func IsPasswordCorrect(username, password string) bool {
	var user User
	db.Where("username = ? and password = ?", username, password).First(&username)
	return user.ID != 0
}
