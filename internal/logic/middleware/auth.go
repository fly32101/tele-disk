package middleware

import (
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	logicAuth "tele-disk/internal/logic/auth"
)

// Auth returns a middleware that validates Bearer token and injects user info into ctx.
func Auth() ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		// Public allow-list: file proxy 无需登录
		if strings.HasPrefix(r.URL.Path, "/api/files/proxy") {
			r.Middleware.Next()
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			r.Response.WriteStatus(http.StatusUnauthorized, g.Map{"error": "missing bearer token"})
			r.Exit()
			return
		}
		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if token == "" {
			r.Response.WriteStatus(http.StatusUnauthorized, g.Map{"error": "empty token"})
			r.Exit()
			return
		}
		claims, err := logicAuth.ParseToken(token)
		if err != nil {
			r.Response.WriteStatus(http.StatusUnauthorized, g.Map{"error": err.Error()})
			r.Exit()
			return
		}
		r.SetCtxVar("userId", claims.UserID)
		r.SetCtxVar("username", claims.Username)
		r.SetCtxVar("isAdmin", claims.IsAdmin)
		r.Middleware.Next()
	}
}
