package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	jwt.RegisteredClaims
	PreferredUsername string `json:"preferred_username"`
	RealmAccess       struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	ResourceAccess map[string]struct {
		Roles []string `json:"roles"`
	} `json:"resource_access"`
	TenantID string `json:"tenant_id,omitempty"`
}

func BasicAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		if !ok {
			c.Header("WWW-Authenticate", `Basic realm="DMS"`)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Basic authentication required"})
			c.Abort()
			return
		}

		claims, err := ValidateBasicAuth(username, password)
		if err != nil {
			log.Printf("Basic auth validation failed: %v", err)
			c.Header("WWW-Authenticate", `Basic realm="DMS"`)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			c.Abort()
			return
		}

		// Extract tenant ID from claims or use a default mapping
		tenantID := extractTenantID(claims)
		if tenantID == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "No tenant access"})
			c.Abort()
			return
		}

		// Add user context to request
		c.Set("user_id", claims.PreferredUsername)
		c.Set("tenant_id", tenantID)
		c.Set("user_roles", claims.RealmAccess.Roles)
		c.Set("claims", claims)

		c.Next()
	}
}

func extractTenantID(claims *UserClaims) string {
	// If tenant_id is explicitly in claims, use it
	if claims.TenantID != "" {
		return claims.TenantID
	}

	// Check for tenant-specific roles or resource access
	for resource, access := range claims.ResourceAccess {
		if strings.HasPrefix(resource, "tenant-") {
			return strings.TrimPrefix(resource, "tenant-")
		}
		// Check if user has access to specific tenant
		for _, role := range access.Roles {
			if strings.HasPrefix(role, "tenant-") {
				return strings.TrimPrefix(role, "tenant-")
			}
		}
	}

	// Default: use username as tenant (for development/simple setups)
	return claims.PreferredUsername
}

func HasRole(c *gin.Context, role string) bool {
	roles, exists := c.Get("user_roles")
	if !exists {
		return false
	}

	userRoles, ok := roles.([]string)
	if !ok {
		return false
	}

	for _, r := range userRoles {
		if r == role {
			return true
		}
	}
	return false
}
