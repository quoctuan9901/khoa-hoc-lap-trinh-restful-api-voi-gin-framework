package v1dto

type LoginInput struct {
	Email    string `json:"email" binding:"required,email,email_advanced"`
	Password string `json:"password" binding:"required,min=8"`
}
