package middleware

import (
	"gateway/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func JwtAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := auth.ExtractId(c, secret)
		if err != nil {
			log.Err(err).Msg("Invalid token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Set("user_id", id)
		c.Next()
	}
}
