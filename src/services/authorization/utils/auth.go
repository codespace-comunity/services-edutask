package utils

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
				"code": http.StatusUnauthorized,
				"error": "authorization invalid",
			})
			return
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
				"code": http.StatusUnauthorized,
				"error": "authorization invalid",
			})
			return
		}

		token := parts[1]

		claims, err := ValidateJWT(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
				"code": http.StatusUnauthorized,
				"error": "jwt not valid",
			})
			return
		}

		if claims.Type != "access_token" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
				"code": http.StatusUnauthorized,
				"error": "jwt not valid",
			})
			return
		}

		c.Set("user_id", claims.Subject)

		c.Next()
	}
}