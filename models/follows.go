package models

import (
	"github.com/jinzhu/gorm"
)

type Follow struct {
	gorm.Model
	Fan  string `gorm:"size:24"`
	Star string `gorm:"size:24"`
}

func init() {
	// create table if not exists
	if !db.HasTable(&Follow{}) {
		db.CreateTable(&Follow{})
	}
}

func AddFollow(fan, star string) error {
	return db.Create(&Follow{Fan: fan, Star: star}).Error
}

func DeleteFollow(fan, star string) error {
	return db.Delete(&Follow{Fan: fan, Star: star}).Error
}

func GetFollows(star string) ([]string, error) {
	var fans []string
	rows, err := db.Select("fan").Where("star = ?", star).Find(&Follow{}).Rows()
	if nil != err {
		return fans, err
	}
	for rows.Next() {
		var fan string
		rows.Scan(&fan)
		fans = append(fans, fan)
	}
	return fans, nil
}

func GetFollowed(fan string) ([]string, error) {
	var stars []string
	rows, err := db.Select("star").Where("fan = ?", fan).Find(&Follow{}).Rows()
	if nil != err {
		return stars, err
	}
	for rows.Next() {
		var star string
		rows.Scan(&star)
		stars = append(stars, star)
	}
	return stars, nil
}

func IsFollowExists(fan, star string) (bool, error) {
	if err := db.Where("fan = ? and star = ?", fan, star).First(&Follow{}).Error; nil != err {
		return false, err
	}
	return true, nil
}
