package routes

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pbearc/codeinsight/models"
	"github.com/pbearc/codeinsight/services"
)

func RegisterRepoRoutes(router *gin.RouterGroup) {
	// Initialize services
	githubService := services.NewGitHubService()
	repoService := services.NewRepoService(githubService)
	llmService := services.NewLLMService()

	// Generate README endpoint
	router.POST("/generate-readme", func(c *gin.Context) {
		var request models.ReadmeGeneratorRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, models.ApiResponse{
				Success: false,
				Error:   "Invalid request: " + err.Error(),
			})
			return
		}

		// Validate inputs
		if request.RepoURL == "" && request.FolderPath == "" {
			c.JSON(http.StatusBadRequest, models.ApiResponse{
				Success: false,
				Error:   "Either repoUrl or folderPath must be provided",
			})
			return
		}

		// Analyze repository structure
		var repoStructure *models.RepoStructure
		var err error
		
		if request.RepoURL != "" {
			repoStructure, err = repoService.AnalyzeGitHubRepository(request.RepoURL)
		} else {
			repoStructure, err = repoService.AnalyzeLocalDirectory(request.FolderPath)
		}
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ApiResponse{
				Success: false,
				Error:   "Failed to analyze repository: " + err.Error(),
			})
			return
		}

		// Format repository structure for LLM
		formattedStructure := repoService.FormatRepositoryStructure(repoStructure)
		
		// Generate README
		readme, err := llmService.GenerateReadme(formattedStructure)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ApiResponse{
				Success: false,
				Error:   "Failed to generate README: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, models.ApiResponse{
			Success: true,
			Data: models.ReadmeResponse{
				Content: readme,
			},
		})
	})

	// Generate project visualization endpoint
	router.POST("/visualize-project", func(c *gin.Context) {
		var request models.ProjectVisualizationRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, models.ApiResponse{
				Success: false,
				Error:   "Invalid request: " + err.Error(),
			})
			return
		}

		// Validate inputs
		if request.RepoURL == "" && request.FolderPath == "" {
			c.JSON(http.StatusBadRequest, models.ApiResponse{
				Success: false,
				Error:   "Either repoUrl or folderPath must be provided",
			})
			return
		}

		// Analyze repository structure
		var repoStructure *models.RepoStructure
		var err error
		
		if request.RepoURL != "" {
			repoStructure, err = repoService.AnalyzeGitHubRepository(request.RepoURL)
		} else {
			repoStructure, err = repoService.AnalyzeLocalDirectory(request.FolderPath)
		}
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ApiResponse{
				Success: false,
				Error:   "Failed to analyze repository: " + err.Error(),
			})
			return
		}

		// Format repository structure for LLM
		formattedStructure := repoService.FormatRepositoryStructure(repoStructure)
		
		// Generate visualization
		diagram, err := llmService.GenerateProjectVisualization(formattedStructure)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ApiResponse{
				Success: false,
				Error:   "Failed to generate visualization: " + err.Error(),
			})
			return
		}

		// Determine diagram type
		diagramType := "flowchart"
		if strings.HasPrefix(diagram, "classDiagram") {
			diagramType = "class"
		} else if strings.HasPrefix(diagram, "sequenceDiagram") {
			diagramType = "sequence"
		} else if strings.HasPrefix(diagram, "stateDiagram") {
			diagramType = "state"
		}

		c.JSON(http.StatusOK, models.ApiResponse{
			Success: true,
			Data: models.ProjectVisualizationResponse{
				DiagramCode: diagram,
				Type:        diagramType,
			},
		})
	})
}