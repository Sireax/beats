package models

type VerificationStatus struct {
	ID   uint `gorm:"primarykey"`
	Name string
}
