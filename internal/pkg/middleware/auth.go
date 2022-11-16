package middleware

import (
	"context"
	"exercise/internal/app/domain"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func WithAuh() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "unauthorized",
			})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "unauthorized",
			})
			c.Abort()
			return
		}

		auths := strings.Split(authHeader, " ")
		if len(auths) != 2 {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "unauthorized",
			})
			c.Abort()
			return
		}
		var user domain.User
		data, err := user.DecryptJWT(auths[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "unauthorized",
			})
			c.Abort()
			return
		}
		userID := int(data["user_id"].(float64))
		ctxUserID := context.WithValue(c.Request.Context(), "user_id", userID)
		c.Request = c.Request.WithContext(ctxUserID)
		c.Next()
	}
}
