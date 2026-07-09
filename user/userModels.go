package user

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	IsAdmin      bool   `json:"is_admin"`
}

type RegisterLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
