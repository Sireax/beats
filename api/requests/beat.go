package requests

type CreateBeatRequest struct {
	Title   string `json:"title" binding:"required"`
	Photo   string `json:"photo" binding:"required"`
	Link    string `json:"link" binding:"required"`
	GenreID uint   `json:"genre_id" binding:"required"`
}

type PurchaseBeatRequest struct {
	LicenseID uint `json:"license_id" binding:"required"`
}

type EditLicenseRequest struct {
	Price      float32 `json:"price"`
	RentalTime string  `json:"rental_time"`
	IsActive   bool    `json:"is_active"`
}
