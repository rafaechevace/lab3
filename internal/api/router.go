// Package api sets up HTTP routes and groups for the application.
package api

import (
	"github.com/networks-security2526/lab3-base/internal/handlers"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns a gin engine with public and authenticated routes.
// Public routes are accessible without authentication. Protected routes are grouped
// and require the AuthMiddleware to succeed before reaching handlers.
func SetupRouter() *gin.Engine {
	// Create a default Gin engine with Logger and Recovery middleware.
	router := gin.Default()

	// Public endpoints: expose application version and authentication endpoints.
	// These do not require a valid user token.
	router.GET("/version", handlers.GetVersion)
	router.POST("/signup", handlers.Signup)
	router.POST("/login", handlers.Login)

	// Group routes that require authentication.
	// AuthMiddleware validates the request (e.g., token) and aborts when invalid.
	authRequired := router.Group("/")
	authRequired.Use(handlers.AuthMiddleware())
	{
		// Document routes operate on a specific user's document identified by
		// :username and :doc_id. Handlers should validate ownership/permissions.
		// POST and PUT both map to CreateOrUpdateDocument: POST creates, PUT updates
		// — using a single handler simplifies shared logic for create-or-update flows.
		authRequired.GET("/:username/:doc_id", handlers.GetDocument)
		authRequired.POST("/:username/:doc_id", handlers.CreateOrUpdateDocument)
		authRequired.PUT("/:username/:doc_id", handlers.CreateOrUpdateDocument)
		authRequired.DELETE("/:username/:doc_id", handlers.DeleteDocument)

		// Retrieve all documents for a given user.
		authRequired.GET("/:username/_all_docs", handlers.GetAllDocuments)
	}

	return router
}
