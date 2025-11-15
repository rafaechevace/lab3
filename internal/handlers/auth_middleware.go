package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware returns a Gin middleware that validates the Authorization header.
// It ensures a token is present, looks it up in the in-memory store, verifies expiration,
// and places the associated username into the request context for downstream handlers.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read the raw Authorization header.
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No header -> immediate unauthorized response.
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Falta la cabecera de autorización"})
			c.Abort()
			return
		}

		// Expect the header to follow the pattern: "token <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "token" {
			// Malformed header -> unauthorized.
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de autorización inválido"})
			c.Abort()
			return
		}

		token := parts[1]

		// Safely read token metadata from the shared token store.
		tokenMutex.Lock()
		tokenInfo, exists := tokenStore[token]
		tokenMutex.Unlock()

		// Token must exist and not be expired.
		if !exists || time.Now().After(tokenInfo.ExpiresAt) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
			c.Abort()
			return
		}

		// Attach authenticated username to the context for later handlers.
		c.Set("username", tokenInfo.Username)

		// Proceed to the next handler in the chain.
		c.Next()
	}
}
