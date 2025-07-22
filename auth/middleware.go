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

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token is required"})
			c.Abort()
			return
		}

		claims, err := ValidateJWT(tokenString)
		if err != nil {
			log.Printf("JWT validation failed: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
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
