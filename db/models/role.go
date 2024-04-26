package models

const (
	ClientRoleID = 1
	ArtistRoleID = 2
)

type Role struct {
	ID   uint   `gorm:"primarykey" json:"id"`
	Name string `json:"name"`
}

func (r *Role) IsClient() bool {
	return r.ID == ClientRoleID
}

func (r *Role) IsArtist() bool {
	return r.ID == ArtistRoleID
}
