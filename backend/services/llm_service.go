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
	apiKey   string
	apiURL   string
	mockMode bool
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
Generate clear documentation for the following code written in {language} as a Markdown document:
{code}

Include the following sections:

Purpose and functionality overview
Parameters and return values (if applicable)
Usage examples
Key components and their interactions
Dependencies and requirements

The documentation should be suitable for developers who may not be familiar with this code.
`,
	"inline_documentation": `
Generate inline documentation comments for the following {language} code. 
Add detailed comments before functions, classes, and complex code blocks. 
For each function or method, document:
- Purpose
- Parameters
- Return values
- Notable behavior or edge cases

Use the appropriate comment style for {language}:
- For JavaScript/TypeScript/Java: Use /** ... */ for functions/methods and // for inline comments
- For Python: Use """ ... """ docstrings for functions/methods and # for inline comments
- For Go: Use // comments following Go's documentation conventions
- For other languages: Follow standard conventions

Here's the code:
{code}

Return the complete code with your added documentation comments.
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
	"readme_generator": `
Generate a comprehensive README.md file for the following repository. Here's the repository structure and file summaries:

{repo_structure}

Based on this information:

1. Create a complete README.md that includes:
   - Project title and description
   - Installation instructions
   - Usage examples
   - Features overview
   - API documentation (if applicable)
   - Contributing guidelines
   - License information

2. The README should be professional and well-structured with proper Markdown formatting.
3. Include badges if applicable (e.g., build status, license).
4. If you see code examples in the file summaries, include them in the usage section.

Return only the README.md content in proper Markdown format.
`,
	"project_visualization": `
Based on the following repository structure:

{repo_structure}

Generate a hierarchical Mermaid flowchart diagram (TD direction - top down) that visualizes the complete folder structure. Requirements:

1. Represent the repository as a hierarchical tree structure with proper parent-child relationships
2. Show ALL directories and their subdirectories
3. Include files within each directory
4. Use different styles for files vs. directories
5. Use the graph TD layout (top-down diagram)

IMPORTANT: 
- Do NOT include markdown code block syntax (three backtick). Return ONLY the raw Mermaid diagram code.
- Make sure to represent the exact hierarchical structure with proper nesting of subdirectories.
- Use descriptive node IDs to avoid collisions when directories or files have the same name.
- Include ALL directories and their relationships.
`,
}

// NewLLMService creates a new LLM service instance
func NewLLMService() *LLMService {
	apiKey := os.Getenv("GEMINI_API_KEY")
	mockMode := apiKey == "" || os.Getenv("USE_MOCK_LLM") == "true"
	
	if mockMode {
		fmt.Println("WARNING: LLM service running in mock mode. Responses will be simulated.")
	}
	
	return &LLMService{
		apiKey:   apiKey,
		apiURL:   "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-pro:generateContent",
		mockMode: mockMode,
	}
}

// CallGeminiLLM calls the Gemini API with the given prompt
func (s *LLMService) CallGeminiLLM(prompt string) (string, error) {
	// Check if we're in mock mode
	if s.mockMode {
		return s.getMockResponse(prompt), nil
	}
	
	// Check if API key is available
	if s.apiKey == "" {
		return "", fmt.Errorf("Gemini API key is not configured")
	}

	// Log request (without the full prompt for privacy)
	fmt.Printf("Making request to Gemini API (prompt length: %d)\n", len(prompt))

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
			"maxOutputTokens": 8000, // Increased for longer responses
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
	client := &http.Client{
		Timeout: 60 * 1000 * 1000 * 1000, // 60 second timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Gemini API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read API response: %w", err)
	}

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Printf("API Error Response: %s\n", string(body))
		return "", fmt.Errorf("Gemini API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
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
		return "", fmt.Errorf("Failed to parse API response: %w (body: %s)", err, string(body))
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("No response content from Gemini API")
	}

	return response.Candidates[0].Content.Parts[0].Text, nil
}

// getMockResponse generates mock responses for testing without an API key
func (s *LLMService) getMockResponse(prompt string) string {
	// Determine what kind of response to generate based on the prompt
	if strings.Contains(prompt, "inline documentation") {
		return "/* MOCK RESPONSE: This is a mock inline documentation response. */\n\n" +
			"/**\n * A mock documented function\n * @param {string} input - The input value\n * @returns {string} The processed result\n */\nfunction mockFunction(input) {\n  // Process the input\n  return input + ' processed';\n}"
	} else if strings.Contains(prompt, "generate a comprehensive README") {
		return "# Mock Project\n\n## Description\nThis is a mock README generated for testing purposes.\n\n## Installation\n```\nnpm install mock-project\n```\n\n## Usage\n```javascript\nconst mock = require('mock-project');\nmock.doSomething();\n```"
	} else if strings.Contains(prompt, "mermaid diagram") {
		return "graph TD\n  A[Main Module] --> B[Component 1]\n  A --> C[Component 2]\n  B --> D[Utility 1]\n  C --> E[Utility 2]"
	} else if strings.Contains(prompt, "documentation") {
		return "# Mock Documentation\n\n## Overview\nThis is a mock documentation response for testing.\n\n## Functions\n- `function1`: Does something interesting\n- `function2`: Does something else"
	} else {
		return "This is a mock response from the LLM service. You're seeing this because either no API key is configured or mock mode is enabled. In a real scenario, this would be a detailed response based on your prompt."
	}
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
func (s *LLMService) GenerateDocumentation(code, language string, inline bool) (*models.DocumentationResponse, error) {
	// Choose template based on whether inline documentation is requested
	templateKey := "documentation"
	if inline {
		templateKey = "inline_documentation"
	}

	// Format the prompt using the selected template
	prompt := s.FormatPrompt(promptTemplates[templateKey], map[string]string{
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
		Inline:        inline,
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

// GenerateReadme generates a README.md file from repository information
func (s *LLMService) GenerateReadme(repoStructure string) (string, error) {
	// Format the prompt using the readme generator template
	prompt := s.FormatPrompt(promptTemplates["readme_generator"], map[string]string{
		"repo_structure": repoStructure,
	})

	// Call the LLM
	response, err := s.CallGeminiLLM(prompt)
	if err != nil {
		return "", fmt.Errorf("README Generation Error: %v", err)
	}

	return response, nil
}

// GenerateProjectVisualization generates a Mermaid diagram of project structure
func (s *LLMService) GenerateProjectVisualization(repoStructure string) (string, error) {
	// Format the prompt using the project visualization template
	prompt := s.FormatPrompt(promptTemplates["project_visualization"], map[string]string{
		"repo_structure": repoStructure,
	})

	// Call the LLM
	response, err := s.CallGeminiLLM(prompt)
	if err != nil {
		return "", fmt.Errorf("Project Visualization Error: %v", err)
	}

	// Ensure the response is a valid Mermaid diagram
	if !strings.Contains(response, "graph ") && 
	   !strings.Contains(response, "flowchart ") && 
	   !strings.Contains(response, "classDiagram") && 
	   !strings.Contains(response, "sequenceDiagram") {
		return "graph TD\n  A[Error] --> B[Invalid diagram generated]\n  B --> C[Please try again]", nil
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