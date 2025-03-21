package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pbearc/codeinsight/routes"
)

func main() {
	// Load environment variables from .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found or error loading it. Using environment variables.")
	}
	
	// Set default port if not specified
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	
	// Create a new Gin router
	router := gin.Default()
	
	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	
	// Global error handler middleware
	router.Use(func(c *gin.Context) {
		c.Next()

		// Handle any panics
		if len(c.Errors) > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Server Error",
				"error":   c.Errors.Last().Error(),
			})
		}
	})

	// API routes group
	api := router.Group("/api")
	
	// Register routes
	routes.RegisterAnalysisRoutes(api)
	routes.RegisterGitHubRoutes(api.Group("/github"))
	routes.RegisterImplementationRoutes(api.Group("/implementations"))
	routes.RegisterRepoRoutes(api.Group("/repo"))
	
	// Register new developer analysis routes
	routes.RegisterDeveloperRoutes(api.Group("/developers"))
	
	// Health check endpoint
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	
	// Start server
	log.Printf("Server starting on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}