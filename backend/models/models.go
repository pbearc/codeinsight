package models

import "encoding/json"

// CodeAnalysisRequest represents a request to analyze code
type CodeAnalysisRequest struct {
	Code         string `json:"code" binding:"required"`
	Language     string `json:"language" default:"javascript"`
	AnalysisType string `json:"analysisType" default:"analyze"`
}

// DocumentationRequest represents a request to generate documentation
type DocumentationRequest struct {
	Code        string `json:"code" binding:"required"`
	Language    string `json:"language" default:"javascript"`
	UseExamples bool   `json:"useExamples" default:"false"`
}

// ImplementationCompareRequest represents a request to compare implementations
type ImplementationCompareRequest struct {
	Implementations []map[string]interface{} `json:"implementations" binding:"required,min=2"`
	Language        string                   `json:"language" default:"javascript"`
}

// Repository represents a GitHub repository
type Repository struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
	Owner       struct {
		Login string `json:"login"`
		URL   string `json:"url"`
	} `json:"owner"`
	Stars int `json:"stars,omitempty"`
}

// File represents a file in a GitHub repository
type File struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	URL     string `json:"url"`
	Content string `json:"content"`
}

// SearchResult represents a GitHub code search result
type SearchResult struct {
	Repository Repository `json:"repository"`
	Path       string     `json:"path"`
	Name       string     `json:"name"`
	URL        string     `json:"url"`
	Score      float64    `json:"score"`
	SHA        string     `json:"sha"`
}

// SearchResponse represents a response from the GitHub search API
type SearchResponse struct {
	TotalCount int           `json:"total_count"`
	Items      []SearchResult `json:"items"`
}

// AnalysisResponse represents an analysis result from the LLM
type AnalysisResponse struct {
	Analysis  string `json:"analysis"`
	Type      string `json:"type"`
	Language  string `json:"language"`
}

// DocumentationResponse represents a documentation result from the LLM
type DocumentationResponse struct {
	Documentation string `json:"documentation"`
	Language      string `json:"language"`
}

// ComparisonResponse represents a comparison result from the LLM
type ComparisonResponse struct {
	Comparison          string `json:"comparison"`
	ImplementationCount int    `json:"implementationCount"`
}

// ApiResponse is a generic API response
type ApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}