package dto

type UserResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Token        string `json:"token"`         //added token for mobile applications
	RefreshToken string `json:"refresh_token"` //added refresh token for mobile applications
}
