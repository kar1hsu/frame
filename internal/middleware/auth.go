package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/karlhsu/frame/internal/pkg/cache"
	jwtpkg "github.com/karlhsu/frame/internal/pkg/jwt"
	"github.com/karlhsu/frame/internal/pkg/response"
)

const (
	CtxUserIDKey   = "user_id"
	CtxUsernameKey = "username"
	CtxRoleCodeKey = "role_code"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "缺少 Authorization 头")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Authorization 格式错误")
			return
		}

		tokenString := parts[1]
		if cache.IsTokenBlacklisted(tokenString) {
			response.Unauthorized(c, "Token 已失效，请重新登录")
			return
		}

		claims, err := jwtpkg.ParseToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "Token 无效或已过期")
			return
		}

		c.Set(CtxUserIDKey, claims.UserID)
		c.Set(CtxUsernameKey, claims.Username)
		c.Set(CtxRoleCodeKey, claims.RoleCode)
		c.Next()
	}
}

func GetUserID(c *gin.Context) uint {
	val, exists := c.Get(CtxUserIDKey)
	if !exists {
		return 0
	}
	return val.(uint)
}

func GetRoleCode(c *gin.Context) string {
	val, exists := c.Get(CtxRoleCodeKey)
	if !exists {
		return ""
	}
	return val.(string)
}
