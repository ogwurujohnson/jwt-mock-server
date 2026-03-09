package token

import (
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// VerifyToken checks the signature and returns the claims if valid
func VerifyToken(tokenString string, pubKey *rsa.PublicKey) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Always validate the signing method is what you expect (RS256)
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return pubKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}