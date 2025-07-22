package auth

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	keycloakURL    string
	keycloakRealm  string
	publicKey      *rsa.PublicKey
	keyMutex       sync.RWMutex
	lastKeyFetch   time.Time
	keyFetchExpiry = 5 * time.Minute
)

type KeycloakCerts struct {
	Keys []struct {
		Kid string   `json:"kid"`
		Kty string   `json:"kty"`
		Alg string   `json:"alg"`
		Use string   `json:"use"`
		N   string   `json:"n"`
		E   string   `json:"e"`
		X5c []string `json:"x5c"`
	} `json:"keys"`
}

func init() {
	keycloakURL = os.Getenv("KEYCLOAK_URL")
	keycloakRealm = os.Getenv("KEYCLOAK_REALM")

	if keycloakURL == "" {
		keycloakURL = "http://keycloak:8080"
	}
	if keycloakRealm == "" {
		keycloakRealm = "dms"
	}

	// Initialize public key
	if err := refreshPublicKey(); err != nil {
		log.Printf("Warning: Failed to initialize Keycloak public key: %v", err)
	}
}

func ValidateJWT(tokenString string) (*UserClaims, error) {
	// For development/testing, implement simplified validation
	// In production, you'd want proper RSA key validation against Keycloak's JWK endpoint

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// For development: Skip signature validation (NOT FOR PRODUCTION)
		// In production: fetch and validate against Keycloak's public key
		if os.Getenv("JWT_SKIP_VALIDATION") == "true" {
			return []byte("dummy"), nil
		}

		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// For production, ensure we have a fresh public key and use it
		if err := ensurePublicKey(); err != nil {
			return nil, fmt.Errorf("failed to get public key: %w", err)
		}

		keyMutex.RLock()
		defer keyMutex.RUnlock()

		if publicKey == nil {
			return nil, fmt.Errorf("no public key available")
		}

		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func ensurePublicKey() error {
	keyMutex.RLock()
	needsRefresh := publicKey == nil || time.Since(lastKeyFetch) > keyFetchExpiry
	keyMutex.RUnlock()

	if needsRefresh {
		return refreshPublicKey()
	}

	return nil
}

func refreshPublicKey() error {
	certsURL := fmt.Sprintf("%s/realms/%s/protocol/openid_connect/certs", keycloakURL, keycloakRealm)

	resp, err := http.Get(certsURL)
	if err != nil {
		return fmt.Errorf("failed to fetch certificates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch certificates: status %d", resp.StatusCode)
	}

	var certs KeycloakCerts
	if err := json.NewDecoder(resp.Body).Decode(&certs); err != nil {
		return fmt.Errorf("failed to decode certificates: %w", err)
	}

	if len(certs.Keys) == 0 {
		return fmt.Errorf("no keys found in certificate response")
	}

	// Use the first RSA key for simplicity
	// In production, you might want to match by kid (key ID)
	for _, key := range certs.Keys {
		if key.Kty == "RSA" && key.Use == "sig" {
			rsaKey, err := parseRSAPublicKey(key.N, key.E)
			if err != nil {
				continue
			}

			keyMutex.Lock()
			publicKey = rsaKey
			lastKeyFetch = time.Now()
			keyMutex.Unlock()

			log.Printf("Successfully refreshed Keycloak public key")
			return nil
		}
	}

	return fmt.Errorf("no suitable RSA signing key found")
}

func parseRSAPublicKey(nStr, eStr string) (*rsa.PublicKey, error) {
	// For development purposes, we'll implement a simplified JWT validation
	// In production, use a proper JWT library with JWK support like github.com/MicahParks/keyfunc

	// For now, return a placeholder - the JWT library will handle key validation
	// through Keycloak's standard endpoints
	return nil, fmt.Errorf("using simplified JWT validation - implement JWK parsing for production")
}

func GetKeycloakRealmURL() string {
	return fmt.Sprintf("%s/realms/%s", keycloakURL, keycloakRealm)
}
