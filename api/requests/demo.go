package requests

type CreateDemoRequest struct {
	BeatID uint   `json:"beat_id" binding:"required"`
	Start  string `json:"start" binding:"required"`
	End    string `json:"end" binding:"required"`
}
