package models

type Verification struct {
	ID                   uint `gorm:"primarykey"`
	BeatID               uint
	Beat                 *Beat `json:"beat"`
	UserID               uint
	User                 *User `json:"user"`
	VerificationStatusID uint
	VerificationStatus   *VerificationStatus `json:"verificationStatus"`
}
