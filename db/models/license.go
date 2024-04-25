package models

type License struct {
	ID            uint    `gorm:"primarykey" json:"id"`
	Price         float32 `json:"price"`
	RentalTime    string  `json:"rentalTime"`
	LicenseTypeID uint
	LicenseType   *LicenseType `json:"licenseType"`
	BeatID        uint
	Beat          *Beat `json:"beat"`
	IsActive      bool  `json:"isActive"`
}
