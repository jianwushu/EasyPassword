package v1

import (
	"easy-password-backend/config"
	"easy-password-backend/internal/apierror"
	"easy-password-backend/internal/crypto"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 创建一个用于 JWT 身份验证的 Gin 中间件。
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			handleError(c, apierror.ErrUnauthorized)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			handleError(c, apierror.New(401, "Authorization header format must be Bearer {token}"))
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := crypto.ValidateJWT(tokenString, cfg.JWTSecret)
		if err != nil {
			handleError(c, apierror.ErrInvalidToken)
			c.Abort()
			return
		}

		// 在上下文中设置用户 ID，以供下游处理程序使用
		c.Set("userID", claims.UserID)

		c.Next()
	}
}