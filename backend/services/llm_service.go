package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/pbearc/codeinsight/models"
)

// LLMService provides methods for interacting with the LLM
type LLMService struct {
	apiKey  string
	apiURL  string
}

// Prompt templates for different analysis tasks
var promptTemplates = map[string]string{
	"analyze": `
Analyze the following code written in {language}:

{language}
{code}

Provide a comprehensive analysis including:

Overview of what the code does
Design patterns used (if any)
Potential bugs or issues
Performance considerations
Suggestions for improvement

Your analysis should be detailed but concise, focusing on the most important aspects.
`,
	"documentation": `
Generate clear documentation for the following code written in {language}:
{code}

Include the following sections:

Purpose and functionality overview
Parameters and return values (if applicable)
Usage examples
Key components and their interactions
Dependencies and requirements

The documentation should be in markdown format and suitable for developers who may not be familiar with this code.
`,
	"explain": `
Explain the following code written in {language} in simple terms:
{code}

Focus on:

What the code does in plain English
How it works, using simple explanations
The purpose of each major component
How the parts work together

Your explanation should be accessible to someone with basic programming knowledge.
`,
	"patternIdentification": `
Identify design patterns in the following code written in {language}:

{code}

For each pattern you identify:

Name the pattern
Explain how it's implemented in this code
Describe the benefits this pattern provides
Suggest any alternative patterns that could also work

Be specific and reference the actual code implementation.
`,
}

// NewLLMService creates a new LLM service instance
func NewLLMService() *LLMService {
	apiKey := os.Getenv("GEMINI_API_KEY")
	return &LLMService{
		apiKey: apiKey,
		apiURL: "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent",
	}
}

// CallGeminiLLM calls the Gemini API with the given prompt
func (s *LLMService) CallGeminiLLM(prompt string) (string, error) {
	// Prepare request payload
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{
						"text": prompt,
					},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":     0.7,
			"maxOutputTokens": 1000,
		},
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Create request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s?key=%s", s.apiURL, s.apiKey), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("Gemini API error: %s", body)
	}

	// Parse response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Gemini API")
	}

	return response.Candidates[0].Content.Parts[0].Text, nil
}

// FormatPrompt formats a prompt using the template and parameters
func (s *LLMService) FormatPrompt(template string, params map[string]string) string {
	prompt := template
	for key, value := range params {
		prompt = strings.ReplaceAll(prompt, "{"+key+"}", value)
	}
	return prompt
}

// AnalyzeCode analyzes code using the Gemini LLM
// AnalyzeCode analyzes code using the Gemini LLM
func (s *LLMService) AnalyzeCode(code, language, analysisType string) (*models.AnalysisResponse, error) {
	// Get the appropriate prompt template
	template, ok := promptTemplates[analysisType]
	if !ok {
		template = promptTemplates["analyze"]
	}

	// Format the prompt
	prompt := s.FormatPrompt(template, map[string]string{
		"code":     code,
		"language": language,
	})

	// Call the LLM
	response, err := s.CallGeminiLLM(prompt)
	if err != nil {
		return nil, fmt.Errorf("Code Analysis Error: %v", err)
	}

	return &models.AnalysisResponse{
		Analysis: response,
		Type:     analysisType,
		Language: language,
	}, nil
}

// GenerateDocumentation generates documentation for code
func (s *LLMService) GenerateDocumentation(code, language string) (*models.DocumentationResponse, error) {
	// Format the prompt using the documentation template
	prompt := s.FormatPrompt(promptTemplates["documentation"], map[string]string{
		"code":     code,
		"language": language,
	})

	// Call the LLM
	response, err := s.CallGeminiLLM(prompt)
	if err != nil {
		return nil, fmt.Errorf("Documentation Generation Error: %v", err)
	}

	return &models.DocumentationResponse{
		Documentation: response,
		Language:      language,
	}, nil
}

// EnhanceWithExamples enhances documentation with examples from GitHub
func (s *LLMService) EnhanceWithExamples(documentation string, examples []map[string]interface{}) (string, error) {
	if examples == nil || len(examples) == 0 {
		return documentation, nil
	}

	// Create a prompt to enhance documentation with real examples
	var examplesText []string
	for _, ex := range examples {
		repoInfo, ok := ex["repository"].(map[string]interface{})
		if !ok {
			continue
		}
		
		fileInfo, ok := ex["file"].(map[string]interface{})
		if !ok {
			continue
		}
		
		fullName, ok := repoInfo["full_name"].(string)
		if !ok {
			continue
		}
		
		content, ok := fileInfo["content"].(string)
		if !ok {
			continue
		}
		
		// Truncate content if too long
		if len(content) > 500 {
			content = content[:500] + "..."
		}
		
		examplesText = append(examplesText, fmt.Sprintf("Example from %s:\n```\n%s\n```", fullName, content))
	}

	prompt := fmt.Sprintf(`
Here is some documentation:
%s
Enhance this documentation with these real-world usage examples:
%s
Create an improved version of the documentation that incorporates insights from these examples.
Make sure to add a "Real-world Examples" section and refer to specific patterns found in the examples.
	`, documentation, strings.Join(examplesText, "\n\n"))

	// Call the LLM
	response, err := s.CallGeminiLLM(prompt)
	if err != nil {
		fmt.Printf("Documentation Enhancement Error: %v\n", err)
		// Return original documentation if enhancement fails
		return documentation, nil
	}

	return response, nil
}

// CompareImplementations compares different implementations of an algorithm or function
func (s *LLMService) CompareImplementations(implementations []map[string]interface{}) (*models.ComparisonResponse, error) {
	if implementations == nil || len(implementations) < 2 {
		return nil, fmt.Errorf("at least two implementations are required for comparison")
	}

	// Create a prompt to compare implementations
	var implementationsText []string
	for i, impl := range implementations {
		repoInfo, ok := impl["repository"].(map[string]interface{})
		if !ok {
			continue
		}
		
		fileInfo, ok := impl["file"].(map[string]interface{})
		if !ok {
			continue
		}
		
		fullName, ok := repoInfo["full_name"].(string)
		if !ok {
			continue
		}
		
		content, ok := fileInfo["content"].(string)
		if !ok {
			continue
		}
		
		// Truncate content if too long
		if len(content) > 1000 {
			content = content[:1000] + "..."
		}
		
		implementationsText = append(implementationsText, fmt.Sprintf("Implementation %d from %s:\n```\n%s\n```", i+1, fullName, content))
	}

	prompt := fmt.Sprintf(`
Compare the following implementations of the same functionality:
%s
Provide a comprehensive comparison including:

Differences in approach
Performance implications
Code quality and readability
Best practices used
Strengths and weaknesses of each approach

Conclude with recommendations on which implementation is best for different scenarios.
	`, strings.Join(implementationsText, "\n\n"))

	// Call the LLM
	response, err := s.CallGeminiLLM(prompt)
	if err != nil {
		return nil, fmt.Errorf("Implementation Comparison Error: %v", err)
	}

	return &models.ComparisonResponse{
		Comparison:          response,
		ImplementationCount: len(implementations),
	}, nil
}