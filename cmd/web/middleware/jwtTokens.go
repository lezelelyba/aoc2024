// Package handles JWT tokens
package middleware

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Generates JWT token.
// Accepts provider name, token, secret to use for encoding and validity period.
// Returns string representation of the Token.
func GenerateJWT(provider, token string, jwtSecret []byte, validityPeriod time.Duration) (string, error) {
	validUntil := fmt.Sprint(time.Now().Add(validityPeriod).Unix())

	claims := jwt.MapClaims{
		"provider":    provider,
		"token":       token,
		"valid_until": fmt.Sprint(validUntil),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString(jwtSecret)
}

// Parses JWT token
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

// Checks if token is valid
func TokenValid(token *jwt.Token) bool {

	// empty token
	if token == nil {
		return false
	}

	// token not parsable by jwt module
	if !token.Valid {
		return false
	}

	// token without claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	// get token expiration
	validUntilStr := claims["valid_until"].(string)

	i, err := strconv.ParseInt(validUntilStr, 10, 64)

	if err != nil {
		return false
	}

	validUntil := time.Unix(i, 0)

	// token is valid if it's not expired
	return !time.Now().After(validUntil)
}
