package models

const (
	ClientRoleType = "client"
	ArtistRoleType = "artist"
)

type Role struct {
	ID   uint   `gorm:"primarykey" json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (r *Role) IsClient() bool {
	return r.Type == ClientRoleType
}

func (r *Role) IsArtist() bool {
	return r.Type == ArtistRoleType
}
