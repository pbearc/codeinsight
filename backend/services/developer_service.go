package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pbearc/codeinsight/models"
)

// DeveloperAnalysisService provides methods for analyzing GitHub developers
type DeveloperAnalysisService struct {
	githubService *GitHubService
	llmService    *LLMService
}

// NewDeveloperAnalysisService creates a new developer analysis service
func NewDeveloperAnalysisService(githubService *GitHubService, llmService *LLMService) *DeveloperAnalysisService {
	return &DeveloperAnalysisService{
		githubService: githubService,
		llmService:    llmService,
	}
}

// AnalyzeDeveloper performs a comprehensive analysis of a GitHub developer profile
func (s *DeveloperAnalysisService) AnalyzeDeveloper(username string, depth string) (*models.DeveloperProfile, error) {
	// Fetch GitHub data for the user
	repoData, err := s.githubService.GetUserRepositories(username)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repositories: %w", err)
	}
	
	// Fetch contribution data
	contributionData, err := s.githubService.GetUserContributions(username)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch contributions: %w", err)
	}
	
	// Fetch language data
	languageData, err := s.githubService.GetUserLanguages(username)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch language data: %w", err)
	}
	
	// Process depth parameter to adjust the amount of data sent to LLM
	var processedRepoData string
	var processedContributionData string
	var processedLanguageData string
	
	switch depth {
	case "light":
		// For light analysis, limit to top 5 repositories
		if len(repoData) > 5 {
			repoData = repoData[:5]
		}
		processedRepoData = s.processRepoDataLight(repoData)
		processedContributionData = s.processContributionDataLight(contributionData)
		processedLanguageData = s.processLanguageDataLight(languageData)
	case "full":
		// For full analysis, use all available data with detailed processing
		processedRepoData = s.processRepoDataFull(repoData)
		processedContributionData = s.processContributionDataFull(contributionData)
		processedLanguageData = s.processLanguageDataFull(languageData)
	default: // medium (default)
		// For medium analysis, use a balanced approach
		if len(repoData) > 10 {
			repoData = repoData[:10]
		}
		processedRepoData = s.processRepoDataMedium(repoData)
		processedContributionData = s.processContributionDataMedium(contributionData)
		processedLanguageData = s.processLanguageDataMedium(languageData)
	}
	
	// Call the LLM service to analyze the developer
	response, err := s.llmService.AnalyzeDeveloper(
		username,
		processedRepoData,
		processedContributionData,
		processedLanguageData,
	)
	if err != nil {
		return nil, fmt.Errorf("LLM analysis failed: %w", err)
	}
	
	// Clean the response to handle potential markdown/backticks
	cleanResponse := extractJSON(response)
	if cleanResponse == "" {
		return nil, fmt.Errorf("could not extract valid JSON from LLM response")
	}
	
	// Parse LLM response
	var profile models.DeveloperProfile
	err = json.Unmarshal([]byte(cleanResponse), &profile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}
	
	return &profile, nil
}

// CompareDevelopers performs a comparison analysis between multiple GitHub developers
func (s *DeveloperAnalysisService) CompareDevelopers(usernames []string, focus string) (*models.ComparativeAnalysis, error) {
	if len(usernames) < 2 {
		return nil, fmt.Errorf("at least two usernames are required for comparison")
	}
	
	// Fetch profiles for all developers
	var developerProfiles []models.DeveloperProfile
	var developerProfilesText string
	
	for _, username := range usernames {
		profile, err := s.AnalyzeDeveloper(username, "medium")
		if err != nil {
			return nil, fmt.Errorf("failed to analyze developer %s: %w", username, err)
		}
		developerProfiles = append(developerProfiles, *profile)
		
		// Convert profile to JSON for prompt
		profileJSON, err := json.Marshal(profile)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal profile for %s: %w", username, err)
		}
		
		developerProfilesText += fmt.Sprintf("--- DEVELOPER: %s ---\n%s\n\n", username, string(profileJSON))
	}
	
	// Call the LLM service to compare the developers
	response, err := s.llmService.CompareDevelopers(developerProfilesText, focus)
	if err != nil {
		return nil, fmt.Errorf("LLM comparison failed: %w", err)
	}
	
	// Clean the response to handle potential markdown/backticks
	cleanResponse := extractJSON(response)
	if cleanResponse == "" {
		return nil, fmt.Errorf("could not extract valid JSON from LLM comparison response")
	}
	
	// Parse LLM response
	var comparison models.ComparativeAnalysis
	err = json.Unmarshal([]byte(cleanResponse), &comparison)
	if err != nil {
		return nil, fmt.Errorf("failed to parse LLM comparison response: %w", err)
	}
	
	return &comparison, nil
}

// extractJSON extracts valid JSON content from a potentially markdown-formatted response
func extractJSON(input string) string {
	// Remove markdown code blocks if present
	// Check for ```json and ``` patterns
	input = removeMarkdownCodeBlocks(input)
	
	// Find the first opening brace
	startIndex := strings.Index(input, "{")
	if startIndex == -1 {
		return ""
	}
	
	// Find the last closing brace
	endIndex := strings.LastIndex(input, "}")
	if endIndex == -1 || endIndex < startIndex {
		return ""
	}
	
	// Extract just the JSON part
	jsonStr := input[startIndex : endIndex+1]
	
	// Validate this is valid JSON by attempting to parse it
	var testMap map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &testMap); err != nil {
		// If parsing fails, try to fix common issues
		jsonStr = cleanJSONString(jsonStr)
		
		// Test again after cleaning
		if err := json.Unmarshal([]byte(jsonStr), &testMap); err != nil {
			return ""
		}
	}
	
	return jsonStr
}

// removeMarkdownCodeBlocks removes markdown code block formatting
func removeMarkdownCodeBlocks(input string) string {
	// Check for ```json or ``` patterns
	lines := strings.Split(input, "\n")
	var result []string
	insideCodeBlock := false
	
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		
		// Check for code block markers
		if strings.HasPrefix(trimmedLine, "```") {
			insideCodeBlock = !insideCodeBlock
			continue // Skip the marker line
		}
		
		// Only include lines that are not code block markers
		if !insideCodeBlock || (insideCodeBlock && trimmedLine != "```") {
			result = append(result, line)
		}
	}
	
	return strings.Join(result, "\n")
}

// cleanJSONString attempts to fix common JSON issues
func cleanJSONString(input string) string {
	// Remove any trailing commas before closing brackets or braces
	regex1 := strings.NewReplacer(
		",]", "]",
		", ]", "]",
		",}", "}",
		", }", "}",
	)
	
	// Replace any invalid escape sequences
	regex2 := strings.NewReplacer(
		"\\'", "'",
		"\\`", "`",
	)
	
	result := regex1.Replace(input)
	result = regex2.Replace(result)
	
	return result
}

// Helper methods to process GitHub data for different analysis depths

func (s *DeveloperAnalysisService) processRepoDataLight(repos []map[string]interface{}) string {
	// For light analysis, just include basic repo info
	lightRepos := make([]map[string]interface{}, 0, len(repos))
	
	for _, repo := range repos {
		lightRepo := map[string]interface{}{
			"name":        repo["name"],
			"description": repo["description"],
			"size":        repo["size"],
			"stars":       repo["stargazers_count"],
			"forks":       repo["forks_count"],
			"language":    repo["language"],
		}
		lightRepos = append(lightRepos, lightRepo)
	}
	
	repoJSON, _ := json.MarshalIndent(lightRepos, "", "  ")
	return string(repoJSON)
}

func (s *DeveloperAnalysisService) processRepoDataMedium(repos []map[string]interface{}) string {
	// For medium analysis, include more repo details
	mediumRepos := make([]map[string]interface{}, 0, len(repos))
	
	for _, repo := range repos {
		languages, _ := repo["languages"].(map[string]interface{})
		
		mediumRepo := map[string]interface{}{
			"name":        repo["name"],
			"description": repo["description"],
			"size":        repo["size"],
			"stars":       repo["stargazers_count"],
			"forks":       repo["forks_count"],
			"language":    repo["language"],
			"created_at":  repo["created_at"],
			"updated_at":  repo["updated_at"],
			"languages":   languages,
		}
		
		// Include commit stats if available
		if commitStats, ok := repo["commit_stats"].(map[string]interface{}); ok {
			mediumRepo["commit_stats"] = map[string]interface{}{
				"total_commits":  commitStats["total_commits"],
				"weekly_average": commitStats["weekly_average"],
			}
		}
		
		mediumRepos = append(mediumRepos, mediumRepo)
	}
	
	repoJSON, _ := json.MarshalIndent(mediumRepos, "", "  ")
	return string(repoJSON)
}

func (s *DeveloperAnalysisService) processRepoDataFull(repos []map[string]interface{}) string {
	// For full analysis, include all repo details
	// Just use the full repo data with indentation for better readability
	repoJSON, _ := json.MarshalIndent(repos, "", "  ")
	return string(repoJSON)
}

func (s *DeveloperAnalysisService) processContributionDataLight(contributions map[string]interface{}) string {
	// For light analysis, just include basic contribution info
	lightContributions := map[string]interface{}{
		"total_events":     contributions["total_events"],
		"monthly_activity": contributions["monthly_activity"],
	}
	
	contribJSON, _ := json.MarshalIndent(lightContributions, "", "  ")
	return string(contribJSON)
}

func (s *DeveloperAnalysisService) processContributionDataMedium(contributions map[string]interface{}) string {
	// For medium analysis, include more contribution details
	// Include all data but limit the events detail if needed
	contribJSON, _ := json.MarshalIndent(contributions, "", "  ")
	return string(contribJSON)
}

func (s *DeveloperAnalysisService) processContributionDataFull(contributions map[string]interface{}) string {
	// For full analysis, include all contribution details
	contribJSON, _ := json.MarshalIndent(contributions, "", "  ")
	return string(contribJSON)
}

func (s *DeveloperAnalysisService) processLanguageDataLight(languages map[string]interface{}) string {
	// For light analysis, just include language percentages
	lightLanguages := map[string]interface{}{
		"percentages": languages["percentages"],
	}
	
	langJSON, _ := json.MarshalIndent(lightLanguages, "", "  ")
	return string(langJSON)
}

func (s *DeveloperAnalysisService) processLanguageDataMedium(languages map[string]interface{}) string {
	// For medium analysis, include more language details
	mediumLanguages := map[string]interface{}{
		"percentages": languages["percentages"],
		"totals":      languages["totals"],
	}
	
	langJSON, _ := json.MarshalIndent(mediumLanguages, "", "  ")
	return string(langJSON)
}

func (s *DeveloperAnalysisService) processLanguageDataFull(languages map[string]interface{}) string {
	// For full analysis, include all language details
	langJSON, _ := json.MarshalIndent(languages, "", "  ")
	return string(langJSON)
}