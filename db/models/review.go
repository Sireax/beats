package models

import (
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	UserID  uint
	User    *User `json:"user"`
	BeatID  uint
	Beat    *Beat  `json:"beat"`
	Content string `json:"content"`
}
