package handler

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestVerify_Token(t *testing.T) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	kid := "test-kid"
	issuer := "http://test"

	h := TokenHandler(privKey, kid, issuer)
	req := httptest.NewRequest(http.MethodGet, "/token", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("GET /token: status = %d; want %d", rec.Code, http.StatusOK)
	}
	body := strings.TrimSpace(rec.Body.String())
	if body == "" {
		t.Error("GET /token: body is empty")
	}
	// JWT has three base64 segments
	parts := strings.Split(body, ".")
	if len(parts) != 3 {
		t.Errorf("GET /token: body does not look like a JWT (got %d parts)", len(parts))
	}
}

func TestVerify_Verify(t *testing.T) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	kid := "test-kid"
	issuer := "http://test"

	tokenHandler := TokenHandler(privKey, kid, issuer)
	tokenReq := httptest.NewRequest(http.MethodGet, "/token", nil)
	tokenRec := httptest.NewRecorder()
	tokenHandler.ServeHTTP(tokenRec, tokenReq)
	if tokenRec.Code != http.StatusOK {
		t.Fatalf("token endpoint: status = %d", tokenRec.Code)
	}
	bearer := strings.TrimSpace(tokenRec.Body.String())
	if bearer == "" {
		t.Fatal("token endpoint: empty body")
	}

	verifyHandler := VerifyHandler(&privKey.PublicKey)
	verifyReq := httptest.NewRequest(http.MethodGet, "/verify", nil)
	verifyReq.Header.Set("Authorization", "Bearer "+bearer)
	verifyRec := httptest.NewRecorder()
	verifyHandler.ServeHTTP(verifyRec, verifyReq)

	if verifyRec.Code != http.StatusOK {
		t.Errorf("GET /verify: status = %d; want %d\nbody: %s", verifyRec.Code, http.StatusOK, verifyRec.Body.String())
	}
	var claims map[string]interface{}
	if err := json.NewDecoder(verifyRec.Body).Decode(&claims); err != nil {
		t.Errorf("GET /verify: invalid JSON: %v", err)
	}
	if claims["iss"] != issuer {
		t.Errorf("GET /verify: claims.iss = %v; want %q", claims["iss"], issuer)
	}
}