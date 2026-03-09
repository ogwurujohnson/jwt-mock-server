package handler

import (
	"encoding/json"
	"github.com/ogwurujohnson/mock-jwk-server/internal/jwk"
	"net/http"
)

func JWKSHandler(keys jwk.JWKS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// RFC 7517 strictly suggests this content type
		w.Header().Set("Content-Type", "application/jwk-set+json")
		json.NewEncoder(w).Encode(keys)
	}
}