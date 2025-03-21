package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pbearc/codeinsight/models"
	"github.com/pbearc/codeinsight/services"
)

// RegisterDeveloperRoutes registers routes for developer analysis
func RegisterDeveloperRoutes(router *gin.RouterGroup) {
	// Initialize services
	githubService := services.NewGitHubService()
	llmService := services.NewLLMService()
	developerService := services.NewDeveloperAnalysisService(githubService, llmService)
	
	// Route for analyzing a single developer
	router.POST("/analyze", func(c *gin.Context) {
		var request models.DeveloperAnalysisRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, models.ApiResponse{
				Success: false,
				Error:   "Invalid request: " + err.Error(),
			})
			return
		}
		
		// Set default depth if not provided
		if request.Depth == "" {
			request.Depth = "medium"
		}
		
		// Analyze the developer
		profile, err := developerService.AnalyzeDeveloper(request.Username, request.Depth)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ApiResponse{
				Success: false,
				Error:   "Developer analysis failed: " + err.Error(),
			})
			return
		}
		
		// Return the analysis
		c.JSON(http.StatusOK, models.ApiResponse{
			Success: true,
			Data: models.DeveloperAnalysisResponse{
				Success: true,
				Profile: *profile,
			},
		})
	})
	
	// Route for comparing multiple developers
	router.POST("/compare", func(c *gin.Context) {
		var request models.DeveloperCompareRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, models.ApiResponse{
				Success: false,
				Error:   "Invalid request: " + err.Error(),
			})
			return
		}
		
		// Validate we have at least 2 usernames
		if len(request.Usernames) < 2 {
			c.JSON(http.StatusBadRequest, models.ApiResponse{
				Success: false,
				Error:   "At least two usernames are required for comparison",
			})
			return
		}
		
		// Set default focus if not provided
		if request.Focus == "" {
			request.Focus = "general"
		}
		
		// Compare the developers
		analysis, err := developerService.CompareDevelopers(request.Usernames, request.Focus)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ApiResponse{
				Success: false,
				Error:   "Developer comparison failed: " + err.Error(),
			})
			return
		}
		
		// Return the comparison
		c.JSON(http.StatusOK, models.ApiResponse{
			Success: true,
			Data: models.DeveloperCompareResponse{
				Success:  true,
				Analysis: *analysis,
			},
		})
	})
}