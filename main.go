package main

import (
	"custom_rugs/auth"
	"custom_rugs/db"
	"custom_rugs/handlers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	db.InitDB("./custom_rugs.db")
	defer db.CloseDB()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("cant load env")
	}
	auth.SetJWTSecret(([]byte)(os.Getenv("JWT_SECRET")))
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.POST("/rug-request", handlers.SubmitRugRequest)
	router.POST("/new-admin", handlers.AddAdminUser)
	router.POST("/login", handlers.Login)
	protectedAPIS := router.Group("/admin")
	protectedAPIS.Use(auth.AuthMiddleware())
	protectedAPIS.GET("/rug-requests", handlers.GetAllRugRequests)
	protectedAPIS.GET("/rug-request/:id", handlers.UpdateRugRequestStatus)

	port := ":8080"
	log.Printf("Server running on port %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)

	}
}
