package handler

import (
	"crypto/rsa"
	"encoding/json"
	"github.com/ogwurujohnson/mock-jwk-server/internal/token"
	"net/http"
	"strings"
)

func VerifyHandler(pubKey *rsa.PublicKey) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Extract token from "Authorization: Bearer <token>"
		authHeader := r.Header.Get("Authorization")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		// 2. Verify
		claims, err := token.VerifyToken(parts[1], pubKey)
		if err != nil {
			http.Error(w, "Token invalid: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// 3. Return claims as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(claims)
	}
}