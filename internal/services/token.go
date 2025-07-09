package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	key []byte
}

func NewAuthService(jwtKey string) *AuthService {
	if jwtKey == "" {
		jwtKey = "very-strong-key"
	}
	return &AuthService{
		[]byte(jwtKey),
	}
}

type Claims struct {
	UserID string `json:"user"`
	jwt.RegisteredClaims
}

func (t *AuthService) GenerateJWT(userId string) (string, error) {
	claims := Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.key)
}

func (t *AuthService) ValidateJWT(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return t.key, nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
