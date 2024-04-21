package models

type LicenseType struct {
	ID   uint   `gorm:"primarykey"`
	Name string `json:"name"`
}
