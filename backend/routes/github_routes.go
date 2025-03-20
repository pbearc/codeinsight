package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pbearc/codeinsight/models"
	"github.com/pbearc/codeinsight/services"
)

func RegisterGitHubRoutes(router *gin.RouterGroup) {
	// Initialize the GitHub service
	githubService := services.NewGitHubService()
	
	// Add request logging middleware
	router.Use(func(c *gin.Context) {
		fmt.Printf("[API Request] %s %s\n", c.Request.Method, c.Request.URL.Path)
		c.Next()
		fmt.Printf("[API Response] Status: %d\n", c.Writer.Status())
	})

	// Search code endpoint
	router.GET("/search", func(c *gin.Context) {
		query := c.Query("query")
		if query == "" {
			c.JSON(http.StatusBadRequest, models.ApiResponse{
				Success: false,
				Error:   "Query is required",
			})
			return
		}

		language := c.Query("language")
		sort := c.DefaultQuery("sort", "best-match")
		order := c.DefaultQuery("order", "desc")
		perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

		// Search code
		results, err := githubService.SearchCode(query, language, sort, order, perPage, page)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ApiResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, models.ApiResponse{
			Success: true,
			Data:    results,
		})
	})

	// Get file content endpoint
	router.GET("/content", func(c *gin.Context) {
		owner := c.Query("owner")
		repo := c.Query("repo")
		path := c.Query("path")

		if owner == "" || repo == "" || path == "" {
			c.JSON(http.StatusBadRequest, models.ApiResponse{
				Success: false,
				Error:   "Owner, repo, and path are required",
			})
			return
		}

		// Get file content
		content, err := githubService.GetFileContent(owner, repo, path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ApiResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, models.ApiResponse{
			Success: true,
			Data:    content,
		})
	})

	// Find library usage examples endpoint
	router.GET("/library-usage", func(c *gin.Context) {
		library := c.Query("library")
		if library == "" {
			c.JSON(http.StatusBadRequest, models.ApiResponse{
				Success: false,
				Error:   "Library name is required",
			})
			return
		}

		language := c.Query("language")
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))

		// Find library usage examples
		examples, err := githubService.FindLibraryUsage(library, language, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ApiResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, models.ApiResponse{
			Success: true,
			Data:    examples,
		})
	})

	// Find implementations endpoint
	router.GET("/implementations", func(c *gin.Context) {
		functionName := c.Query("functionName")
		if functionName == "" {
			c.JSON(http.StatusBadRequest, models.ApiResponse{
				Success: false,
				Error:   "Function name is required",
			})
			return
		}

		language := c.Query("language")
		limit, err := strconv.Atoi(c.DefaultQuery("limit", "3"))
		if err != nil || limit < 1 {
			limit = 3
		} else if limit > 10 {
			limit = 10 // Set a reasonable maximum
		}

		// Use GitHub service with panic recovery
		var implementations []map[string]interface{}
		var implErr error
		
		func() {
			// Recover from any panics
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Recovered from panic: %v\n", r)
					implErr = fmt.Errorf("internal server error")
				}
			}()
			
			// Execute the GitHub service call
			implementations, implErr = githubService.FindImplementations(functionName, language, limit)
		}()

		// Handle any errors
		if implErr != nil {
			fmt.Printf("Error finding implementations: %v\n", implErr)
			c.JSON(http.StatusInternalServerError, models.ApiResponse{
				Success: false,
				Error:   implErr.Error(),
			})
			return
		}

		// Success response
		c.JSON(http.StatusOK, models.ApiResponse{
			Success: true,
			Data:    implementations,
		})
	})
}