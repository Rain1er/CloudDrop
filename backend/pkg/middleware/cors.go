package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles Cross-Origin Resource Sharing for local React dev servers.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if isAllowedLocalOrigin(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		}
		c.Header("Vary", "Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, Content-Length, X-Requested-With, Accept, Origin")
		c.Header("Access-Control-Expose-Headers", "Content-Length")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func isAllowedLocalOrigin(origin string) bool {
	if origin == "" {
		return false
	}
	allowedPrefixes := []string{
		"http://localhost:",
		"http://127.0.0.1:",
		"http://[::1]:",
	}
	for _, prefix := range allowedPrefixes {
		if strings.HasPrefix(origin, prefix) {
			return true
		}
	}
	return false
}
