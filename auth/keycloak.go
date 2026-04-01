package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	keycloakURL          string
	keycloakRealm        string
	keycloakClientID     string
	keycloakClientSecret string
	publicKey            *rsa.PublicKey
	keyMutex             sync.RWMutex
	lastKeyFetch         time.Time
	keyFetchExpiry       = 5 * time.Minute
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
	keycloakClientID = os.Getenv("KEYCLOAK_CLIENT_ID")
	keycloakClientSecret = os.Getenv("KEYCLOAK_CLIENT_SECRET")

	if keycloakURL == "" {
		keycloakURL = "http://keycloak:8082"
	}
	if keycloakRealm == "" {
		keycloakRealm = "dms"
	}
	if keycloakClientID == "" {
		keycloakClientID = "dms-service"
	}
	if keycloakClientSecret == "" {
		keycloakClientSecret = "your-service-secret-key"
	}

	// Initialize public key
	if err := refreshPublicKey(); err != nil {
		log.Printf("Warning: Failed to initialize Keycloak public key: %v", err)
	}
}

func ValidateJWT(tokenString string) (*UserClaims, error) {
	if os.Getenv("JWT_SKIP_VALIDATION") == "true" {
		parser := jwt.NewParser(jwt.WithoutClaimsValidation())
		token, _, err := parser.ParseUnverified(tokenString, &UserClaims{})
		if err != nil {
			return nil, fmt.Errorf("failed to parse token: %w", err)
		}
		claims, ok := token.Claims.(*UserClaims)
		if !ok {
			return nil, fmt.Errorf("invalid token claims")
		}
		return claims, nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

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
	certsURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", keycloakURL, keycloakRealm)

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
	nBytes, err := base64.RawURLEncoding.DecodeString(nStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %w", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(eStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %w", err)
	}

	n := new(big.Int).SetBytes(nBytes)
	e := int(new(big.Int).SetBytes(eBytes).Int64())

	return &rsa.PublicKey{N: n, E: e}, nil
}

func ValidateBasicAuth(username, password string) (*UserClaims, error) {
	// Validate credentials against Keycloak's token endpoint
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", keycloakURL, keycloakRealm)

	data := fmt.Sprintf("grant_type=password&client_id=%s&client_secret=%s&username=%s&password=%s",
		keycloakClientID, keycloakClientSecret, username, password)

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to validate credentials: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid credentials")
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	// Parse the JWT token to extract claims
	return ValidateJWT(tokenResponse.AccessToken)
}

func GetKeycloakRealmURL() string {
	return fmt.Sprintf("%s/realms/%s", keycloakURL, keycloakRealm)
}

func GetTokenForUser(username, password string) (string, error) {
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", keycloakURL, keycloakRealm)

	data := fmt.Sprintf("grant_type=password&client_id=%s&client_secret=%s&username=%s&password=%s",
		keycloakClientID, keycloakClientSecret, username, password)

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get token: status %d", resp.StatusCode)
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	return tokenResponse.AccessToken, nil
}
