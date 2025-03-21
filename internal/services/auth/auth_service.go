package service

import (
	"github.com/golang-jwt/jwt/v4"
	"log"
	"time"
)

type (
	AuthService interface {
		GenerateToken(userId string) (*Tokens, error)
		ValidateToken(token string) (string, error)
		RefreshToken(token string) (*Tokens, error)
	}

	authService struct {
		jwtSecret        []byte
		jwtRefreshSecret []byte
	}

	Tokens struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	claims struct {
		UserID string `json:"user_id"`
		jwt.RegisteredClaims
	}

	refreshClaims struct {
		UserID string `json:"user_id"`
		jwt.RegisteredClaims
	}
)

func NewAuthService() AuthService {
	return &authService{}
}

func (au *authService) GenerateToken(userId string) (*Tokens, error) {
	//create jwt token
	expirationTime := time.Now().Add(72 * time.Hour) //72 hours expiration time
	claims := &claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, tokenErr := token.SignedString(au.jwtSecret)
	if tokenErr != nil {
		log.Print(tokenErr)
		return nil, tokenErr
	}
	//create jwt refresh token
	refreshExpirationTime := time.Now().Add(7 * 24 * time.Hour) //7 days
	rClaims := refreshClaims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rClaims)
	refreshTokenString, refreshTokenErr := refreshToken.SignedString(au.jwtRefreshSecret)
	if refreshTokenErr != nil {
		log.Print(refreshTokenErr)
		return nil, refreshTokenErr
	}

	return &Tokens{
		AccessToken:  tokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (au *authService) ValidateToken(token string) (string, error) {

}

func (au *authService) RefreshToken(token string) (*Tokens, error) {

}
