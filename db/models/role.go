package models

type Role struct {
	ID   uint   `gorm:"primarykey" json:"id"`
	Name string `json:"name"`
}
