package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pbearc/codeinsight/models"
	"github.com/pbearc/codeinsight/services"
)

func RegisterAnalysisRoutes(router *gin.RouterGroup) {
	// Initialize the LLM service
	llmService := services.NewLLMService()

	// Analyze code endpoint
	router.POST("/analyze", func(c *gin.Context) {
		var request models.CodeAnalysisRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, models.ApiResponse{
				Success: false,
				Error:   "Invalid request: " + err.Error(),
			})
			return
		}

		if request.Code == "" {
			c.JSON(http.StatusBadRequest, models.ApiResponse{
				Success: false,
				Error:   "Code is required",
			})
			return
		}

		// Set default values if not provided
		if request.Language == "" {
			request.Language = "javascript"
		}
		if request.AnalysisType == "" {
			request.AnalysisType = "analyze"
		}

		// Analyze the code
		result, err := llmService.AnalyzeCode(request.Code, request.Language, request.AnalysisType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ApiResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, models.ApiResponse{
			Success: true,
			Data:    result,
		})
	})

	// Generate documentation endpoint
	router.POST("/generate-docs", func(c *gin.Context) {
		var request models.DocumentationRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, models.ApiResponse{
				Success: false,
				Error:   "Invalid request: " + err.Error(),
			})
			return
		}

		if request.Code == "" {
			c.JSON(http.StatusBadRequest, models.ApiResponse{
				Success: false,
				Error:   "Code is required",
			})
			return
		}

		// Set default values if not provided
		if request.Language == "" {
			request.Language = "javascript"
		}

		// Generate documentation
		result, err := llmService.GenerateDocumentation(request.Code, request.Language)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ApiResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, models.ApiResponse{
			Success: true,
			Data:    result,
		})
	})

	// Compare implementations endpoint
	router.POST("/implementations/compare", func(c *gin.Context) {
		var request models.ImplementationCompareRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, models.ApiResponse{
				Success: false,
				Error:   "Invalid request: " + err.Error(),
			})
			return
		}

		if len(request.Implementations) < 2 {
			c.JSON(http.StatusBadRequest, models.ApiResponse{
				Success: false,
				Error:   "At least two implementations are required for comparison",
			})
			return
		}

		// Compare implementations
		result, err := llmService.CompareImplementations(request.Implementations)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ApiResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, models.ApiResponse{
			Success: true,
			Data:    result,
		})
	})
}