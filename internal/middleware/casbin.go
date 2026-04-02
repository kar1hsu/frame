package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/karlhsu/frame/internal/app"
	"github.com/karlhsu/frame/internal/pkg/response"
)

func CasbinRBAC() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleCode := GetRoleCode(c)
		if roleCode == "" {
			response.Forbidden(c, "无法获取用户角色")
			return
		}

		// super admin bypasses all checks
		if roleCode == "admin" {
			c.Next()
			return
		}

		obj := c.Request.URL.Path
		act := c.Request.Method

		ok, err := app.Enforcer.Enforce(roleCode, obj, act)
		if err != nil {
			app.Log.Errorw("casbin enforce error", "error", err)
			response.Forbidden(c, "权限校验异常")
			return
		}
		if !ok {
			response.Forbidden(c, "无访问权限")
			return
		}

		c.Next()
	}
}
