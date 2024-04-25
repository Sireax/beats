package models

type LicenseType struct {
	ID   uint   `gorm:"primarykey" json:"id"`
	Name string `json:"name"`
}
