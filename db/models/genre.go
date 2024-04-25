package models

type Genre struct {
	ID   uint   `gorm:"primarykey" json:"id"`
	Name string `json:"name"`
}
