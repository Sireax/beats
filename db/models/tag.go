package models

type Tag struct {
	ID   uint `gorm:"primarykey"`
	Name string
}
