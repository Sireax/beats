package requests

type CreateBeatRequest struct {
	Title   string `json:"title" binding:"required"`
	Photo   string `json:"photo" binding:"required"`
	Link    string `json:"link" binding:"required"`
	GenreID uint   `json:"genre_id" binding:"required"`
}
