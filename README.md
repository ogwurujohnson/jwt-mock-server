# JWT Mock Server

A minimal HTTP server that acts as a **mock OAuth2/OIDC-style issuer**: it exposes a JWKS endpoint, issues RS256 JWTs, and can verify Bearer tokens. Useful for local development, integration tests, or learning JWT/JWKS flows.

## Features

- **JWKS** — Serves a JSON Web Key Set at `/.well-known/jwks.json` (RFC 7517)
- **Token issuance** — `GET /token` returns a signed JWT (RS256, 1-hour expiry)
- **Token verification** — `GET/POST /verify` validates `Authorization: Bearer <token>` and returns the decoded claims as JSON
- **RSA 2048** — Keys are generated at startup; no persistent key storage

## Requirements

- Go 1.26+

## Quick Start

```bash
# Run locally
go run ./cmd/server
```

Server listens on **http://localhost:8080**.

### With Docker

The image is published at [Docker Hub](https://hub.docker.com/r/ogwurujohnson/jwtmockserver):

```bash
docker pull ogwurujohnson/jwtmockserver
docker run -p 8080:8080 ogwurujohnson/jwtmockserver
```

Or build from source:

```bash
docker build -t jwt-mock-server .
docker run -p 8080:8080 jwt-mock-server
```

## API Reference

### `GET /.well-known/jwks.json`

Returns the public keys in JWK Set format. Content-Type: `application/jwk-set+json`.

**Example response:**

```json
{
  "keys": [
    {
      "kty": "RSA",
      "use": "sig",
      "kid": "mock-key-id-1",
      "alg": "RS256",
      "n": "...",
      "e": "AQAB"
    }
  ]
}
```

### `GET /token`

Issues a new JWT signed with the server’s private key.

**Response:** Plain text body containing the JWT string.

**Token claims (example):**

| Claim | Value |
|-------|--------|
| `iss` | `http://localhost:8080` |
| `sub` | `1234567890` |
| `name` | `Mock User` |
| `iat` | Issued-at timestamp |
| `exp` | Expiry (1 hour from issue) |

**Example:**

```bash
curl -s http://localhost:8080/token
```

### `GET /verify` or `POST /verify`

Verifies a JWT from the `Authorization` header and returns the claims as JSON.

**Headers:**

- `Authorization: Bearer <your-jwt>`

**Success (200):** JSON body with the token’s claims.

**Errors:**

- `401` — Missing/invalid `Authorization` header or invalid/expired token.

**Example:**

```bash
TOKEN=$(curl -s http://localhost:8080/token)
curl -s -H "Authorization: Bearer $TOKEN" http://localhost:8080/verify
```

## Project Layout

```
.
├── cmd/server/          # Entrypoint: key generation, routes, HTTP server
├── internal/
│   ├── handler/         # HTTP handlers: JWKS, token, verify
│   ├── jwk/             # RSA public key → JWK (RFC 7517)
│   └── token/           # JWT creation (RS256) and verification
├── Dockerfile
├── go.mod
└── README.md
```

## Development

```bash
# Run tests
go test ./...

# Build binary
go build -o server ./cmd/server
./server
```

## License

See repository for license details.
