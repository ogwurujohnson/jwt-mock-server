package main

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"github.com/ogwurujohnson/mock-jwk-server/internal/handler"
	"github.com/ogwurujohnson/mock-jwk-server/internal/token"
	"github.com/ogwurujohnson/mock-jwk-server/internal/jwk"
	"net/http"
)

func main() {
	// 1. Generate a key pair
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	
	// 2. Convert to RFC 7517 JWK
	key := jwk.PublicKeyToJWK(&privateKey.PublicKey, "mock-key-id-1")
	set := jwk.JWKS{Keys: []jwk.JWK{key}}

	// 3. Setup Routes
	http.HandleFunc("/.well-known/jwks.json", handler.JWKSHandler(set))

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
    // We reuse the same privateKey and kid generated at startup
    tString, err := token.GenerateToken(privateKey, "mock-key-id-1", "http://localhost:8080")
    if err != nil {
        http.Error(w, "Failed to generate token", 500)
        return
    }

    w.Write([]byte(tString))
	})

	log.Println("Mock JWK Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
