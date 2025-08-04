package v1dto

type LoginInput struct {
	Email    string `json:"email" binding:"required,email,email_advanced"`
	Password string `json:"password" binding:"required,min=8"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RequestPasswordInput struct {
	Email string `json:"email" binding:"required,email,email_advanced"`
}

type ResetPasswordInput struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}
