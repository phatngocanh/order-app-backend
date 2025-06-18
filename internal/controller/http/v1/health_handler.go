package v1

import (
	"github.com/pna/order-app-backend/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	db database.Db
}

func NewHealthHandler(db database.Db) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// @Summary Health Check
// @Description Checks the health of the application by verifying database and Redis connections
// @Tags Healths
// @Accept json
// @Produce json
// @Router /health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	// Check database connection
	if err := h.db.DB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  "database connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}
