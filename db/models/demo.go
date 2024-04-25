package models

import "time"

type Demo struct {
	ID    uint      `gorm:"primarykey" json:"id"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}
