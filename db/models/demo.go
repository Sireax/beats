package models

type Demo struct {
	ID     uint   `gorm:"primarykey" json:"id"`
	Start  string `json:"start"`
	End    string `json:"end"`
	BeatID uint
	Beat   *Beat `json:"beat"`
}
