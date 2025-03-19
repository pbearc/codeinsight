package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/pbearc/codeinsight/models"
)

type GitHubService struct {
	token    string
	baseURL  string
	headers  map[string]string
}

// NewGitHubService creates a new GitHub service instance
func NewGitHubService() *GitHubService {
	token := os.Getenv("GITHUB_TOKEN")
	return &GitHubService{
		token:   token,
		baseURL: "https://api.github.com",
		headers: map[string]string{
			"Accept":        "application/vnd.github.v3+json",
			"Authorization": "token " + token,
			"User-Agent":    "CodeInsight-App",
		},
	}
}

// makeRequest makes a request to the GitHub API
func (s *GitHubService) makeRequest(method, endpoint string, params map[string]string) ([]byte, error) {
	// Build URL with query parameters
	url := s.baseURL + endpoint
	if len(params) > 0 && method == "GET" {
		url += "?"
		for key, value := range params {
			url += key + "=" + value + "&"
		}
		url = strings.TrimSuffix(url, "&")
	}

	// Create request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// Add headers
	for key, value := range s.headers {
		req.Header.Add(key, value)
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error: %s", body)
	}

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// SearchCode searches for code in GitHub repositories
func (s *GitHubService) SearchCode(query, language, sort, order string, perPage, page int) (*models.SearchResponse, error) {
	// Build query string
	queryString := query
	if language != "" {
		queryString += " language:" + language
	}

	// Build params
	params := map[string]string{
		"q":        queryString,
		"sort":     sort,
		"order":    order,
		"per_page": fmt.Sprintf("%d", perPage),
		"page":     fmt.Sprintf("%d", page),
	}

	// Make request
	body, err := s.makeRequest("GET", "/search/code", params)
	if err != nil {
		return nil, err
	}

	// Parse response
	var searchResponse struct {
		TotalCount int `json:"total_count"`
		Items      []struct {
			Name       string `json:"name"`
			Path       string `json:"path"`
			SHA        string `json:"sha"`
			URL        string `json:"html_url"`
			Score      float64 `json:"score"`
			Repository struct {
				Name        string `json:"name"`
				FullName    string `json:"full_name"`
				HTMLURL     string `json:"html_url"`
				Description string `json:"description"`
				Owner       struct {
					Login   string `json:"login"`
					HTMLURL string `json:"html_url"`
				} `json:"owner"`
			} `json:"repository"`
		} `json:"items"`
	}
	if err := json.Unmarshal(body, &searchResponse); err != nil {
		return nil, err
	}

	// Process results
	result := &models.SearchResponse{
		TotalCount: searchResponse.TotalCount,
		Items:      make([]models.SearchResult, 0, len(searchResponse.Items)),
	}

	for _, item := range searchResponse.Items {
		searchResult := models.SearchResult{
			Name:  item.Name,
			Path:  item.Path,
			URL:   item.URL,
			Score: item.Score,
			SHA:   item.SHA,
			Repository: models.Repository{
				Name:        item.Repository.Name,
				FullName:    item.Repository.FullName,
				URL:         item.Repository.HTMLURL,
				Description: item.Repository.Description,
				Owner: struct {
					Login string `json:"login"`
					URL   string `json:"url"`
				}{
					Login: item.Repository.Owner.Login,
					URL:   item.Repository.Owner.HTMLURL,
				},
			},
		}
		result.Items = append(result.Items, searchResult)
	}

	return result, nil
}

// GetFileContent gets content of a file from GitHub
func (s *GitHubService) GetFileContent(owner, repo, path string) (string, error) {
	// Make request
	body, err := s.makeRequest("GET", fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, path), nil)
	if err != nil {
		return "", err
	}

	// Parse response
	var fileResponse struct {
		Content string `json:"content"`
	}
	if err := json.Unmarshal(body, &fileResponse); err != nil {
		return "", err
	}

	// Decode base64 content
	decodedContent, err := base64.StdEncoding.DecodeString(fileResponse.Content)
	if err != nil {
		return "", err
	}

	return string(decodedContent), nil
}

// FindLibraryUsage finds library usage examples in GitHub repositories
func (s *GitHubService) FindLibraryUsage(library, language string, limit int) ([]map[string]interface{}, error) {
	// Build query string
	queryString := library
	if language != "" {
		queryString += " language:" + language

		// Add language-specific import patterns
		switch strings.ToLower(language) {
		case "javascript":
			queryString += fmt.Sprintf(" \"require('%s')\" OR \"import %s\" OR \"import { \" OR \"from '%s'\"", library, library, library)
		case "python":
			queryString += fmt.Sprintf(" \"import %s\" OR \"from %s\"", library, library)
		case "java":
			queryString += fmt.Sprintf(" \"import %s\"", library)
		}
	}

	// Filter out test files and documentation
	queryString += " NOT filename:test NOT filename:spec NOT path:test NOT path:docs"

	// Build params
	params := map[string]string{
		"q":        queryString,
		"sort":     "stars",
		"order":    "desc",
		"per_page": fmt.Sprintf("%d", limit),
	}

	// Make request
	body, err := s.makeRequest("GET", "/search/code", params)
	if err != nil {
		return nil, err
	}

	// Parse response
	var searchResponse struct {
		Items []struct {
			Name       string `json:"name"`
			Path       string `json:"path"`
			URL        string `json:"html_url"`
			Repository struct {
				Name        string `json:"name"`
				FullName    string `json:"full_name"`
				HTMLURL     string `json:"html_url"`
				Owner       struct {
					Login string `json:"login"`
				} `json:"owner"`
				StargazersCount int `json:"stargazers_count"`
			} `json:"repository"`
		} `json:"items"`
	}
	if err := json.Unmarshal(body, &searchResponse); err != nil {
		return nil, err
	}

	// Get file content for each result
	examples := make([]map[string]interface{}, 0, limit)
	for _, item := range searchResponse.Items {
		if len(examples) >= limit {
			break
		}

		owner := item.Repository.Owner.Login
		repo := item.Repository.Name
		path := item.Path

		// Get file content
		content, err := s.GetFileContent(owner, repo, path)
		if err != nil {
			fmt.Printf("Error fetching file content for %s: %s\n", path, err)
			continue
		}

		example := map[string]interface{}{
			"repository": map[string]interface{}{
				"name":     item.Repository.Name,
				"full_name": item.Repository.FullName,
				"url":      item.Repository.HTMLURL,
				"stars":    item.Repository.StargazersCount,
			},
			"file": map[string]interface{}{
				"name":    item.Name,
				"path":    item.Path,
				"url":     item.URL,
				"content": content,
			},
		}
		examples = append(examples, example)
	}

	return examples, nil
}

// FindImplementations finds implementations of a specific algorithm or function
func (s *GitHubService) FindImplementations(functionName, language string, limit int) ([]map[string]interface{}, error) {
	// Build query string
	queryString := "function " + functionName
	if language != "" {
		queryString += " language:" + language

		// Add language-specific function patterns
		switch strings.ToLower(language) {
		case "javascript":
			queryString += fmt.Sprintf(" \"function %s\" OR \"const %s = function\" OR \"const %s =\"", functionName, functionName, functionName)
		case "python":
			queryString += fmt.Sprintf(" \"def %s\"", functionName)
		case "java":
			queryString += fmt.Sprintf(" \"public * %s\" OR \"private * %s\"", functionName, functionName)
		}
	}

	// Build params
	params := map[string]string{
		"q":        queryString,
		"sort":     "stars",
		"order":    "desc",
		"per_page": fmt.Sprintf("%d", limit),
	}

	// Make request
	body, err := s.makeRequest("GET", "/search/code", params)
	if err != nil {
		return nil, err
	}

	// Parse response
	var searchResponse struct {
		Items []struct {
			Name       string `json:"name"`
			Path       string `json:"path"`
			URL        string `json:"html_url"`
			Repository struct {
				Name        string `json:"name"`
				FullName    string `json:"full_name"`
				HTMLURL     string `json:"html_url"`
				Owner       struct {
					Login string `json:"login"`
				} `json:"owner"`
			} `json:"repository"`
		} `json:"items"`
	}
	if err := json.Unmarshal(body, &searchResponse); err != nil {
		return nil, err
	}

	// Get file content for each result
	implementations := make([]map[string]interface{}, 0, limit)
	for _, item := range searchResponse.Items {
		if len(implementations) >= limit {
			break
		}

		owner := item.Repository.Owner.Login
		repo := item.Repository.Name
		path := item.Path

		// Get file content
		content, err := s.GetFileContent(owner, repo, path)
		if err != nil {
			fmt.Printf("Error fetching file content for %s: %s\n", path, err)
			continue
		}

		implementation := map[string]interface{}{
			"repository": map[string]interface{}{
				"name":     item.Repository.Name,
				"full_name": item.Repository.FullName,
				"url":      item.Repository.HTMLURL,
			},
			"file": map[string]interface{}{
				"name":    item.Name,
				"path":    item.Path,
				"url":     item.URL,
				"content": content,
			},
		}
		implementations = append(implementations, implementation)
	}

	return implementations, nil
}