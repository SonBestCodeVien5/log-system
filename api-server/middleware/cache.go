package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// DashboardNoCache prevents browsers from keeping an old dashboard shell or
// asset bundle while the bind-mounted static files change during development.
func DashboardNoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/" || strings.HasPrefix(path, "/assets/") {
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}
		c.Next()
	}
}
