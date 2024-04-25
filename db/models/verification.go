package models

type Verification struct {
	ID                   uint `gorm:"primarykey" json:"id"`
	BeatID               uint
	Beat                 *Beat `json:"beat"`
	UserID               uint
	User                 *User `json:"user"`
	VerificationStatusID uint
	VerificationStatus   *VerificationStatus `json:"verificationStatus"`
}
