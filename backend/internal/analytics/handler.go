package analytics

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAnalyticsHandler(c *gin.Context) {
	data, err := GetMachineAnalytics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
