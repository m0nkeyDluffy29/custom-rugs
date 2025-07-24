package main

import (
	"custom_rugs/db"
	"custom_rugs/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB("./custom_rugs.db")
	defer db.CloseDB()

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

	router.GET("/rug-requests", handlers.GetAllRugRequests)
	router.GET("/rug-request/:id", handlers.UpdateRugRequestStatus)

	port := ":8080"
	log.Printf("Server running on port %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)

	}
}
