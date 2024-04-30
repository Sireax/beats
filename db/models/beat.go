package models

import "time"

type Beat struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	ReleaseDate time.Time `json:"releaseDate"`
	Photo       string    `json:"photo"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	GenreID     uint      `json:"genreID"`
	Genre       *Genre    `json:"genre"`
	UserID      uint      `json:"userID"`
	User        *User     `json:"user"`
	IsHide      bool      `json:"isHide"`

	Demos    []*Demo    `json:"demos"`
	Snippets []*Snippet `json:"snippets"`
	Licenses []*License `json:"licenses"`
	Tags     []*Tag     `gorm:"many2many:song_tags;" json:"tags"`
}
