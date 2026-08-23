// package routes

// import (
// 	"net/http"

// 	alerts "factorypulse/backend/internal/alert"
// 	"factorypulse/backend/internal/analytics"
// 	"factorypulse/backend/internal/auth"
// 	"factorypulse/backend/internal/machines"
// 	"factorypulse/backend/internal/maintenance"
// 	"factorypulse/backend/internal/middleware"
// 	"factorypulse/backend/internal/sensors"

// 	"github.com/gin-gonic/gin"
// )

// func Setup(r *gin.Engine) {
// 	r.GET("/health", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{"status": "ok"})
// 	})

// 	r.POST("/api/auth/register", auth.RegisterHandler)
// 	r.POST("/api/auth/login", auth.LoginHandler)

// 	protected := r.Group("/api")
// 	protected.Use(middleware.AuthRequired())
// 	{
// 		protected.GET("/me", func(c *gin.Context) {
// 			userID, _ := c.Get("user_id")
// 			role, _ := c.Get("role")
// 			c.JSON(http.StatusOK, gin.H{"user_id": userID, "role": role})
// 		})

// 		protected.POST("/machines", machines.CreateMachineHandler)
// 		protected.GET("/machines", machines.GetAllMachinesHandler)
// 		protected.GET("/machines/:id", machines.GetMachineHandler)

// 		protected.POST("/readings", sensors.CreateReadingHandler)
// 		protected.GET("/machines/:id/latest-reading", sensors.GetLatestReadingHandler)
// 		protected.GET("/machines/:id/history", sensors.GetReadingHistoryHandler)
// 		protected.GET("/alerts", alerts.GetActiveAlertsHandler)
// 		protected.POST("/maintenance", maintenance.CreateJobHandler)
// 		protected.GET("/maintenance", maintenance.GetAllJobsHandler)
// 		protected.PATCH("/maintenance/:id/status", maintenance.UpdateStatusHandler)
// 		protected.PATCH("/maintenance/:id/assign", maintenance.AssignJobHandler)
// 		protected.GET("/analytics", analytics.GetAnalyticsHandler)
// 	}
// }

package routes

import (
	"net/http"

	// "factorypulse/backend/internal/alerts"
	alerts "factorypulse/backend/internal/alert"
	"factorypulse/backend/internal/analytics"
	"factorypulse/backend/internal/auth"
	"factorypulse/backend/internal/machines"
	"factorypulse/backend/internal/maintenance"
	"factorypulse/backend/internal/middleware"
	"factorypulse/backend/internal/sensors"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/api/auth/register", auth.RegisterHandler)
	r.POST("/api/auth/login", auth.LoginHandler)

	protected := r.Group("/api")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/me", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			role, _ := c.Get("role")
			c.JSON(http.StatusOK, gin.H{"user_id": userID, "role": role})
		})

		protected.GET("/machines", machines.GetAllMachinesHandler)
		protected.GET("/machines/:id", machines.GetMachineHandler)
		protected.POST("/machines", middleware.RoleRequired("admin"), machines.CreateMachineHandler)

		protected.GET("/machines/:id/latest-reading", sensors.GetLatestReadingHandler)
		protected.GET("/machines/:id/history", sensors.GetReadingHistoryHandler)
		protected.POST("/readings", sensors.CreateReadingHandler)

		protected.GET("/alerts", alerts.GetActiveAlertsHandler)

		protected.GET("/maintenance", maintenance.GetAllJobsHandler)
		protected.POST("/maintenance", middleware.RoleRequired("admin", "engineer"), maintenance.CreateJobHandler)
		protected.PATCH("/maintenance/:id/status", maintenance.UpdateStatusHandler)
		protected.PATCH("/maintenance/:id/assign", maintenance.AssignJobHandler)

		protected.GET("/analytics", middleware.RoleRequired("admin", "production_manager"), analytics.GetAnalyticsHandler)
	}
}
