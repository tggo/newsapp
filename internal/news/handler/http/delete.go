package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) delete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}