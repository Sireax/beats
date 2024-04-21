package models

import "time"

type Snippet struct {
	ID    uint      `gorm:"primarykey"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}
