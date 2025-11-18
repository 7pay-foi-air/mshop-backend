package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SECRET_KEY []byte
var REFRESH_SECRET []byte

func SetAccesSecretKey(secret string) {
	SECRET_KEY = []byte(secret)
}

func SetRefreshSecretKey(secret string) {
	REFRESH_SECRET = []byte(secret)
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	OrgID  string `json:"org_id"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID, role, orgID string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		OrgID:  orgID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SECRET_KEY)
}

func GenerateRefreshToken(userID string) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if len(REFRESH_SECRET) > 0 {
		return token.SignedString(REFRESH_SECRET)
	}
	return token.SignedString(SECRET_KEY)
}

func ValidateAccessToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid access token")
	}
	return claims, nil
}

func ValidateRefreshToken(tokenStr string) (string, error) {
	parseWith := func(secret []byte) (string, error) {
		token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil {
			return "", err
		}
		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok || !token.Valid {
			return "", errors.New("invalid refresh token")
		}
		return claims.Subject, nil
	}

	if len(REFRESH_SECRET) > 0 {
		if sub, err := parseWith(REFRESH_SECRET); err == nil {
			return sub, nil
		}
	}

	return parseWith(SECRET_KEY)
}
