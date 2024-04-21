package models

import "gorm.io/gorm"

type Purchase struct {
	gorm.Model
	UserID    uint
	User      *User `json:"user"`
	BeatID    uint
	Beat      *Beat `json:"beat"`
	LicenseID uint
	License   *License `json:"license"`
}
