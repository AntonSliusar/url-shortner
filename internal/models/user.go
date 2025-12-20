package models

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	GoogleID     string `json:"google_id"`
	Role         string `json:"role"`
	IsVerified   bool   `json:"email_verified"`
}
