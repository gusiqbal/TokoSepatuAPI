package middleware

import (
	"learnapirest/helpers"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := helpers.VerifyJWT(tokenStr, secret)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}

		c.Set("userId", userID)
		c.Next()
	}
}
