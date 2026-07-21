package jwt

import (
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	RoleID string `json:"role_id"`
	jwtlib.RegisteredClaims
}

func GenerateToken(userID, roleID, secret string, expiry time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		RoleID: roleID,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenStr, secret string) (*Claims, error) {
	token, err := jwtlib.ParseWithClaims(tokenStr, &Claims{}, func(t *jwtlib.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwtlib.ErrTokenInvalidClaims
}
