package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/CulipBlue/backend_ednic/internal/modules/auth"
	"github.com/CulipBlue/backend_ednic/internal/shared/response"
)

const (
	ContextUserID = "user_id"
	ContextRole   = "role"
	ContextEmail  = "email"
)

func RequireAuth(authService auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenValue := bearerToken(c.GetHeader("Authorization"))
		if tokenValue == "" {
			response.Error(c, http.StatusUnauthorized, "Authentication required", nil)
			c.Abort()
			return
		}

		claims, err := authService.ParseToken(tokenValue)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid or expired token", nil)
			c.Abort()
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextRole, claims.Role)
		c.Set(ContextEmail, claims.Email)
		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get(ContextRole)
		if role != auth.RoleAdmin && role != auth.RoleSuperAdmin {
			response.Error(c, http.StatusForbidden, "Admin access required", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

func CurrentUserID(c *gin.Context) (uint64, bool) {
	value, exists := c.Get(ContextUserID)
	if !exists {
		return 0, false
	}

	switch typed := value.(type) {
	case uint64:
		return typed, true
	case string:
		parsed, err := strconv.ParseUint(typed, 10, 64)
		return parsed, err == nil
	default:
		return 0, false
	}
}

func bearerToken(header string) string {
	header = strings.TrimSpace(header)
	if header == "" {
		return ""
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	return strings.TrimSpace(parts[1])
}
