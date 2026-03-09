package main

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"github.com/ogwurujohnson/mock-jwk-server/internal/handler"
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

	http.HandleFunc("/token", handler.TokenHandler(privateKey, "mock-key-id-1", "http://localhost:8080"))

	http.HandleFunc("/verify", handler.VerifyHandler(&privateKey.PublicKey))

	log.Println("Mock JWK Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
