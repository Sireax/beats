package models

type Tag struct {
	ID   uint `gorm:"primarykey" json:"id"`
	Name string
}
