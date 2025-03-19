package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"your-username/codeinsight/routes"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Global error handler middleware
	r.Use(func(c *gin.Context) {
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

	// Redirect root to Swagger docs
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	// Health check endpoint
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api")
	{
		// Analysis routes
		routes.RegisterAnalysisRoutes(api)

		// GitHub routes
		github := api.Group("/github")
		routes.RegisterGitHubRoutes(github)
	}

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	// Start server
	log.Printf("Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}