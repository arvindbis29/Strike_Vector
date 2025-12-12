package main

import (
	"fmt"
	"net/http"
	handler "voice-hack-backend/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting the server...")

	port := "8080"
	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.New()

	// Add CORS middleware with more permissive settings for development
	ginEngine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins for development
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false, // Set to false when using wildcard origins
	}))

	handler.RouteRequests(ginEngine)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: ginEngine,
	}

	if serverErr := server.ListenAndServe(); serverErr != nil {
		fmt.Println("Some issue occured while initiating the server: " + serverErr.Error())
	}
}
