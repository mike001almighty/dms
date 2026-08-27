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
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := ValidateJWT(tokenString)
		if err != nil {
			log.Printf("JWT validation failed: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		tenantID := extractTenantID(claims)
		if tenantID == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "No tenant access"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.PreferredUsername)
		c.Set("tenant_id", tenantID)
		c.Set("user_roles", claims.RealmAccess.Roles)
		c.Set("claims", claims)

		c.Next()
	}
}

func extractTenantID(claims *UserClaims) string {
	if claims.TenantID != "" {
		return claims.TenantID
	}
	return ""
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

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !HasRole(c, role) {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient role"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireAnyRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, role := range roles {
			if HasRole(c, role) {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient role"})
		c.Abort()
	}
}
