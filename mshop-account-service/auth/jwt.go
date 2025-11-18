package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SECRET_KEY []byte

func SetSecretKey(secret string) {
	if secret == "" {
		panic("JWT_SECRET env var is not set")
	}
	SECRET_KEY = []byte(secret)
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	OrgID  string `json:"org_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID, role, orgID string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		OrgID:  orgID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return SECRET_KEY, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
