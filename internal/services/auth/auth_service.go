package service

import (
	"errors"
	service "github.com/SilentPlaces/basicauth.git/internal/services/vault"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/wire"
	"log"
	"time"
)

const (
	tokenExpireTime        = time.Hour * 72     // 72 hours
	refreshTokenExpireTime = time.Hour * 24 * 7 // 7 days
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

	Claims struct {
		UserID string `json:"user_id"`
		jwt.RegisteredClaims
	}

	RefreshClaims struct {
		UserID string `json:"user_id"`
		jwt.RegisteredClaims
	}
)

func NewAuthService(vault service.SecureVaultService) AuthService {
	jwtConfig, err := vault.GetJWTConfig()
	if err != nil {
		log.Panic(err)
	}
	return &authService{
		jwtSecret:        jwtConfig.JwtSecret,
		jwtRefreshSecret: jwtConfig.JwtRefreshSecret,
	}
}

func (au *authService) GenerateToken(userId string) (*Tokens, error) {
	// Create access token
	expirationTime := time.Now().Add(tokenExpireTime)
	claims := &Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(au.jwtSecret)
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		return nil, err
	}

	// Create refresh token
	refreshExpirationTime := time.Now().Add(refreshTokenExpireTime)
	rClaims := &RefreshClaims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rClaims)
	refreshTokenString, err := refreshToken.SignedString(au.jwtRefreshSecret)
	if err != nil {
		log.Printf("Error generating refresh token: %v", err)
		return nil, err
	}

	return &Tokens{
		AccessToken:  tokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (au *authService) ValidateToken(token string) (string, error) {
	if token == "" {
		log.Print("ValidateToken: token is empty")
		return "", errors.New("token is empty")
	}

	cToken, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return au.jwtSecret, nil
	})
	if err != nil {
		log.Printf("ValidateToken: error parsing token: %v", err)
		return "", err
	}

	claims, ok := cToken.Claims.(*Claims)
	if !ok || !cToken.Valid {
		log.Print("ValidateToken: invalid token")
		return "", errors.New("invalid token")
	}

	return claims.UserID, nil
}

func (au *authService) RefreshToken(token string) (*Tokens, error) {
	if token == "" {
		log.Print("RefreshToken: token is empty")
		return nil, errors.New("token is empty")
	}

	cToken, err := jwt.ParseWithClaims(token, &RefreshClaims{}, func(t *jwt.Token) (interface{}, error) {
		return au.jwtRefreshSecret, nil
	})
	if err != nil {
		log.Printf("RefreshToken: error parsing refresh token: %v", err)
		return nil, err
	}

	refreshClaims, ok := cToken.Claims.(*RefreshClaims)
	if !ok || !cToken.Valid {
		log.Print("RefreshToken: invalid token")
		return nil, errors.New("invalid token")
	}

	return au.GenerateToken(refreshClaims.UserID)
}

var AuthServiceProviderSet = wire.NewSet(NewAuthService)
