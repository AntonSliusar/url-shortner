package models

type VerifyCodeInput struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

type LoginResponse struct {
}