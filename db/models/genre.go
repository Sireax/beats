package models

type Genre struct {
	ID   uint   `gorm:"primarykey"`
	Name string `json:"name"`
}
