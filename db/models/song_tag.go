package models

type SongTag struct {
	ID     uint  `gorm:"primarykey" json:"id"`
	BeatID uint  `json:"beatID"`
	Beat   *Beat `json:"beat"`
	TagID  uint  `json:"tagID"`
	Tag    *Tag  `json:"tag"`
}
