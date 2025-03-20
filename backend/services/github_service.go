package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

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
	
	fmt.Println("Initializing GitHub service...")
	if token == "" {
		fmt.Println("WARNING: GITHUB_TOKEN not set. Using unauthenticated requests (limited to 60 per hour).")
		return &GitHubService{
			token:   "",
			baseURL: "https://api.github.com",
			headers: map[string]string{
				"Accept":     "application/vnd.github.v3+json",
				"User-Agent": "CodeInsight-App",
			},
		}
	}
	
	fmt.Println("Using authenticated GitHub requests with provided token.")
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
    fullURL := s.baseURL + endpoint
    if len(params) > 0 && method == "GET" {
        queryParams := url.Values{}
        for key, value := range params {
            queryParams.Add(key, value)
        }
        fullURL += "?" + queryParams.Encode()
    }

    // Create request
    req, err := http.NewRequest(method, fullURL, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    // Add headers
    for key, value := range s.headers {
        req.Header.Add(key, value)
    }

    // Send request with timeout
    client := &http.Client{
        Timeout: 20 * time.Second,
    }
    
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    // Read response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %w", err)
    }

    // Check rate limit
    if resp.StatusCode == 403 && strings.Contains(string(body), "rate limit") {
        return nil, fmt.Errorf("GitHub API rate limit exceeded. Please wait or use an authenticated token")
    }

    // Check status code
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return nil, fmt.Errorf("GitHub API error (status %d): %s", resp.StatusCode, string(body))
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
    // Validate inputs
    if owner == "" || repo == "" || path == "" {
        return "", fmt.Errorf("owner, repo, and path must all be provided")
    }
    
    // URL encode the path (handle spaces and special characters)
    encodedPath := url.PathEscape(path)
    
    // Make request
    endpoint := fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, encodedPath)
    body, err := s.makeRequest("GET", endpoint, nil)
    if err != nil {
        return "", err
    }
    
    // Parse response
    var fileResponse struct {
        Content  string `json:"content"`
        Encoding string `json:"encoding"`
        Size     int    `json:"size"`
    }
    
    if err := json.Unmarshal(body, &fileResponse); err != nil {
        return "", fmt.Errorf("failed to parse file content response: %w", err)
    }
    
    if fileResponse.Content == "" {
        return "", fmt.Errorf("empty file content returned")
    }
    
    // Decode base64 content
    contentStr := strings.ReplaceAll(fileResponse.Content, "\n", "")
    decodedContent, err := base64.StdEncoding.DecodeString(contentStr)
    if err != nil {
        return "", fmt.Errorf("failed to decode base64 content: %w", err)
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
		"per_page": fmt.Sprintf("%d", limit*2), // Request more results to account for potential filtering
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

	// Process results with concurrency control
	examples := make([]map[string]interface{}, 0, limit)
	var mu sync.Mutex
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 3) // Limit concurrent requests
	
	for _, item := range searchResponse.Items {
		if len(examples) >= limit {
			break
		}

		wg.Add(1)
		go func(item struct {
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
		}) {
			defer wg.Done()
			semaphore <- struct{}{} // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore
			
			owner := item.Repository.Owner.Login
			repo := item.Repository.Name
			path := item.Path

			// Get file content
			content, err := s.GetFileContent(owner, repo, path)
			if err != nil {
				return
			}

			example := map[string]interface{}{
				"repository": map[string]interface{}{
					"name":      item.Repository.Name,
					"full_name": item.Repository.FullName,
					"url":       item.Repository.HTMLURL,
					"stars":     item.Repository.StargazersCount,
				},
				"file": map[string]interface{}{
					"name":    item.Name,
					"path":    item.Path,
					"url":     item.URL,
					"content": content,
				},
			}
			
			mu.Lock()
			if len(examples) < limit {
				examples = append(examples, example)
			}
			mu.Unlock()
		}(item)
	}
	
	wg.Wait()

	return examples, nil
}

// FindImplementations finds implementations of a specific algorithm or function
func (s *GitHubService) FindImplementations(functionName, language string, limit int) ([]map[string]interface{}, error) {
    // Validate input
    if functionName == "" {
        return nil, fmt.Errorf("function name cannot be empty")
    }
    
    // Build query for searching GitHub
    query := s.buildQueryForFunction(functionName, language)
    
    // Search parameters
    params := map[string]string{
        "q":        query,
        "sort":     "stars", // Sort by repository stars to get quality implementations
        "order":    "desc",
        "per_page": fmt.Sprintf("%d", limit * 2), // Request more than needed to account for filtering
    }
    
    // Search GitHub code
    body, err := s.makeRequest("GET", "/search/code", params)
    if err != nil {
        return nil, fmt.Errorf("GitHub search failed: %w", err)
    }
    
    // Parse response
    var searchResponse struct {
        TotalCount int `json:"total_count"`
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
        return nil, fmt.Errorf("failed to parse GitHub response: %w", err)
    }
    
    // No results found
    if len(searchResponse.Items) == 0 {
        return []map[string]interface{}{}, nil
    }
    
    // Process results with concurrency, but limit to avoid rate limits
    var mu sync.Mutex // Mutex to protect the shared results slice
    
    implementations := make([]map[string]interface{}, 0, limit)
    itemCount := len(searchResponse.Items)
    if itemCount > limit*2 {
        itemCount = limit*2 // Limit the number of items to process
    }
    
    // Use a semaphore to limit concurrent API calls
    semaphore := make(chan struct{}, 3) // At most 3 concurrent requests
    var wg sync.WaitGroup
    
    // Process each search result
    for _, item := range searchResponse.Items[:itemCount] {
        wg.Add(1)
        go func(item struct {
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
        }) {
            defer wg.Done()
            
            // Acquire semaphore slot
            semaphore <- struct{}{}
            defer func() { <-semaphore }()
            
            owner := item.Repository.Owner.Login
            repo := item.Repository.Name
            path := item.Path
            
            // Get file content
            content, err := s.GetFileContent(owner, repo, path)
            if err != nil {
                return
            }
            
            // Create implementation data
            implementation := map[string]interface{}{
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
            
            // Add to results if we haven't reached the limit
            mu.Lock()
            if len(implementations) < limit {
                implementations = append(implementations, implementation)
            }
            mu.Unlock()
        }(item)
    }
    
    // Wait for all goroutines to complete
    wg.Wait()
    
    return implementations, nil
}

// buildQueryForFunction builds an optimized query string for finding function implementations
func (s *GitHubService) buildQueryForFunction(functionName, language string) string {
    // Handle special characters in function name
    escapedName := strings.Replace(functionName, " ", "_", -1)
    
    // Base query with language filter
    query := escapedName
    
    // Add language-specific patterns
    if language != "" {
        query += fmt.Sprintf(" language:%s", language)
        
        switch strings.ToLower(language) {
        case "javascript":
            // For JavaScript, look for function declarations, arrow functions, etc.
            functionPatterns := []string{
                fmt.Sprintf("function %s", escapedName),
                fmt.Sprintf("const %s =", escapedName),
                fmt.Sprintf("let %s =", escapedName),
                fmt.Sprintf("var %s =", escapedName),
            }
            query += " (" + strings.Join(functionPatterns, " OR ") + ")"
            
        case "python":
            // For Python, look for function definitions
            query += fmt.Sprintf(" def %s", escapedName)
            
        case "java":
            // For Java, look for method declarations
            methodPatterns := []string{
                fmt.Sprintf("public static .* %s", escapedName),
                fmt.Sprintf("private static .* %s", escapedName),
                fmt.Sprintf("public .* %s", escapedName),
                fmt.Sprintf("private .* %s", escapedName),
            }
            query += " (" + strings.Join(methodPatterns, " OR ") + ")"
            
        case "go":
            // For Go, look for function declarations
            query += fmt.Sprintf(" func %s", escapedName)
            
        // Add other languages as needed
        }
    }
    
    // Exclude test files and documentation
    query += " -filename:test -filename:spec -path:test"
    
    return query
}