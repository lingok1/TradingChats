package api

import (
	"net/http"
	"strings"
	"trading-chats-backend/internal/config"
	"trading-chats-backend/internal/models"
	"trading-chats-backend/internal/service"

	"github.com/gin-gonic/gin"
)

const authContextKey = "auth_context"

func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse(401, "missing bearer token"))
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		authCtx, err := authService.ValidateAccessToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse(401, err.Error()))
			return
		}

		c.Set(authContextKey, authCtx)
		c.Request = c.Request.WithContext(models.WithAuthContext(c.Request.Context(), authCtx))
		c.Next()
	}
}

func RequireRoles(roles ...models.Role) gin.HandlerFunc {
	allowed := map[models.Role]struct{}{}
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(c *gin.Context) {
		authCtx := MustGetAuthContext(c)
		if _, ok := allowed[authCtx.Role]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, models.ErrorResponse(403, "permission denied"))
			return
		}
		c.Next()
	}
}

func MustGetAuthContext(c *gin.Context) *models.AuthContext {
	value, ok := c.Get(authContextKey)
	if !ok {
		return models.GetAuthContext(c.Request.Context())
	}
	authCtx, _ := value.(*models.AuthContext)
	if authCtx == nil {
		return models.GetAuthContext(c.Request.Context())
	}
	return authCtx
}

func SwaggerBasicAuth(cfg *config.SwaggerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.AllowWithoutAuth {
			c.Next()
			return
		}

		user, pass, ok := c.Request.BasicAuth()
		if !ok || user != cfg.Username || pass != cfg.Password {
			c.Header("WWW-Authenticate", `Basic realm="Swagger API Documentation"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Unauthorized: Valid authentication credentials required for Swagger documentation",
			})
			return
		}

		c.Next()
	}
}
