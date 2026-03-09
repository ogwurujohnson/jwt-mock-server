package handler

import (
	"crypto/rsa"
	"encoding/json"
	"net/http"

	"github.com/ogwurujohnson/mock-jwk-server/internal/jwk"
	"github.com/ogwurujohnson/mock-jwk-server/internal/token"
)

func JWKSHandler(keys jwk.JWKS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// RFC 7517 strictly suggests this content type
		w.Header().Set("Content-Type", "application/jwk-set+json")
		json.NewEncoder(w).Encode(keys)
	}
}

// TokenHandler returns an HTTP handler that issues a JWT for the given key and kid.
func TokenHandler(privKey *rsa.PrivateKey, kid, issuer string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tString, err := token.GenerateToken(privKey, kid, issuer)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(tString))
	}
}
