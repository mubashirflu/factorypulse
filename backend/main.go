// package main

// import (
// 	"factorypulse/backend/internal/database"
// 	routes "factorypulse/backend/internal/router"
// 	"log"
// 	"os"

// 	"github.com/gin-gonic/gin"
// 	"github.com/joho/godotenv"
// )

// func main() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Println("warning!.env file not found")
// 	}
// 	database.Connect()
// 	defer database.Pool.Close()
// 	router := gin.Default()

// 	routes.Setup(router)
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "8080"
// 	}
// 	router.Run(":" + port)
// }

package main

import (
	"factorypulse/backend/internal/database"
	routes "factorypulse/backend/internal/router"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	database.Connect()
	defer database.Pool.Close()

	r := gin.Default()

	// CORS enable karo — frontend (5173) ko backend (8080) se baat karne do
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.Setup(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
