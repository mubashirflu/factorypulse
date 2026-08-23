// package alerts

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

//	func GetActiveAlertsHandler(c *gin.Context) {
//		list, err := GetActiveAlerts()
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//			return
//		}
//		c.JSON(http.StatusOK, list)
//	}
package alerts

import (
	"net/http"
	"strconv"

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

func GetAllAlertsHandler(c *gin.Context) {
	list, err := GetAllAlerts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func ResolveAlertHandler(c *gin.Context) {
	alertID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert id"})
		return
	}

	if err := ResolveAlert(alertID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "alert resolved"})
}
