package handlers
import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetVersion returns the current version of the Database as a Service API.
func GetVersion(c *gin.Context) {
	c.String(http.StatusOK, "Database as a Service v0.1.0")
}
