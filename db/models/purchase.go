package models

import (
	"time"
)

type Purchase struct {
	ID        uint `gorm:"primarykey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uint
	User      *User `json:"user"`
	BeatID    uint
	Beat      *Beat `json:"beat"`
	LicenseID uint
	License   *License `json:"license"`
	Purchased int      `json:"purchased" gorm:"-"`
}
