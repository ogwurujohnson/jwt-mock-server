package token

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(privKey *rsa.PrivateKey, kid string, issuer string) (string, error) {
	// 1. Define the Payload (Claims)
	claims := jwt.MapClaims{
		"iss": issuer,                               // Must match your discovery URL
		"sub": "1234567890",                         // The "User ID"
		"name": "Mock User",
		"iat": time.Now().Unix(),                    // Issued At
		"exp": time.Now().Add(time.Hour * 1).Unix(), // Expiration (1 hour)
	}

	// 2. Create the token object with RS256
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// 3. IMPORTANT: Set the 'kid' in the header so the client knows which key to use
	token.Header["kid"] = kid

	// 4. Sign the token with the PRIVATE key
	return token.SignedString(privKey)
}