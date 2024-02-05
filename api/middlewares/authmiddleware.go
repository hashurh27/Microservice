package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("Username")
		password := c.GetHeader("Password")

		// Check if the provided username and password match the expected values
		if username == "admin" && password == "admin" {
			// Authentication successful, proceed to the next middleware or handler
			c.Next()
		} else {
			// Authentication failed, respond with Unauthorized status
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort() // Abort further processing
		}
	}
}
