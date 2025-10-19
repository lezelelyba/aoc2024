package middleware

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(provider, token, secret string, validityPeriod time.Duration) (string, error) {
	validUntil := fmt.Sprint(time.Now().Add(validityPeriod).Unix())

	claims := jwt.MapClaims{
		"provider":    provider,
		"token":       token,
		"valid_until": fmt.Sprint(validUntil),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString([]byte(secret))
}

func ParseToken(tokenStr string, jwtSecret []byte) (*jwt.Token, error) {
	if tokenStr == "" {
		return nil, fmt.Errorf("token is empty")
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})

	return token, err
}

func TokenValid(token *jwt.Token) bool {

	if token == nil {
		return false
	}

	if !token.Valid {
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	validUntilStr := claims["valid_until"].(string)

	i, err := strconv.ParseInt(validUntilStr, 10, 64)

	if err != nil {
		return false
	}

	validUntil := time.Unix(i, 0)

	return !time.Now().After(validUntil)
}
