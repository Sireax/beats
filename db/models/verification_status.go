package models

type VerificationStatus struct {
	ID   uint `gorm:"primarykey" json:"id"`
	Name string
}
