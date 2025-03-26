package refresh_token

type RefreshTokenResDTO struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
