package middleware

import (
	"net/http"
	"rate-limiter/pkg/util"

	"github.com/gin-gonic/gin"
)

func AuthAPI() gin.HandlerFunc {
	apiKey := util.GetEnv("API_KEY", "fallback")
	return func(c *gin.Context) {
		inputKey := c.Request.Header["Api-Key"]
		clientID := c.Request.Header["X-Client-Id"]

		if len(clientID) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"code":    401,
				"message": "Unauthorized Client ID",
			})
		}

		if len(inputKey) == 0 || inputKey[0] != apiKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"code":    401,
				"message": "Unauthorized API Key",
			})
		}

		c.Set("clientID", clientID[0])

		c.Next()
	}
}
