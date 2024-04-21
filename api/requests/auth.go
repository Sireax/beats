package requests

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Photo    string `json:"photo"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
