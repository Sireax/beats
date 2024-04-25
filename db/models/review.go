package models

import "time"

type Review struct {
	ID        uint `gorm:"primarykey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uint
	User      *User `json:"user"`
	BeatID    uint
	Beat      *Beat  `json:"beat"`
	Content   string `json:"content"`
}
