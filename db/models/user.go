package models

type User struct {
	ID       uint   `gorm:"primarykey" json:"id"`
	Username string `json:"name"`
	Email    string `json:"email"`
	Password []byte `json:"-"`
	RoleID   uint   `json:"-"`
	Role     *Role  `json:"role"`
}
