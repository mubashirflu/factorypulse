package alerts

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetActiveAlertsHandler(c *gin.Context) {
	list, err := GetActiveAlerts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}
