package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags HealthCheck
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealtCheckHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Health Check Success",
	})
}
