package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pbearc/codeinsight/models"
	"github.com/pbearc/codeinsight/services"
)

func RegisterImplementationRoutes(router *gin.RouterGroup) {
	// Create a new comparison service
	comparisonService := services.NewComparisonService()
	
	// Compare implementations endpoint
	router.POST("/compare", func(c *gin.Context) {
		// Parse request body
		var request models.ImplementationCompareRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, models.ApiResponse{
				Success: false,
				Error:   "Invalid request body: " + err.Error(),
			})
			return
		}
		
		// Validate request data
		if len(request.Implementations) < 2 {
			c.JSON(http.StatusBadRequest, models.ApiResponse{
				Success: false,
				Error:   "At least 2 implementations required for comparison",
			})
			return
		}
		
		// Generate comparison with panic recovery
		var comparison string
		var err error
		
		func() {
			// Recover from any panics
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Recovered from panic in comparison generation: %v\n", r)
					err = fmt.Errorf("internal server error while generating comparison")
				}
			}()
			
			// Generate the comparison
			comparison, err = comparisonService.CompareImplementations(
				request.Implementations, 
				request.Language,
			)
		}()
		
		// Handle errors
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ApiResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		
		// Return response
		c.JSON(http.StatusOK, models.ApiResponse{
			Success: true,
			Data: map[string]interface{}{
				"comparison": comparison,
			},
		})
	})
}