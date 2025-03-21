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
	}  else if strings.Contains(prompt, "developer profile analysis") || strings.Contains(prompt, "GitHub user profile") {
		return `{"username":"mockuser","executiveSummary":"This is a mock developer profile analysis for testing purposes.","languageAnalysis":[{"name":"JavaScript","percentage":65.5,"lines":12500,"proficiency":"Expert","experience":4.2},{"name":"Python","percentage":25.3,"lines":4800,"proficiency":"Intermediate","experience":2.1}],"projectAnalysis":[{"name":"mock-project-1","size":12500,"stars":42,"forks":15,"complexityRank":"medium","score":7.2}],"contributionPatterns":{"totalCommits":512,"averageCommitsMonth":42.7,"consistencyScore":85.2,"prAcceptanceRate":92.5,"reviewParticipation":75.0,"monthlyActivity":[45,40,38,42,45,50,55,48,42,38,40,45]},"learningVelocity":{"adaptabilityScore":82.5,"technologyAdoption":{"TypeScript":"1 year ago","React":"2 years ago"},"growthTrajectory":"accelerating","complexityTrend":[{"month":1,"score":5.2},{"month":12,"score":7.5}]},"specializations":{"primaryDomains":["Frontend Development","Web Applications"],"secondaryDomains":["API Design","DevOps"],"frameworkExpertise":{"React":85.0,"Node.js":75.0},"topicFocus":{"UI/UX":65.0,"Performance":72.0}},"skillAssessment":{"codeQuality":80.0,"documentation":65.0,"testing":70.0,"performance":85.0,"security":60.0,"collaboration":90.0},"technicalStrengths":["JavaScript Ecosystem","Frontend Architecture","Performance Optimization"],"growthAreas":["Documentation Practices","Security Best Practices"]}`
	} else if strings.Contains(prompt, "developer comparison") || strings.Contains(prompt, "Compare the following GitHub developer profiles") {
		return `{"developers":["mockuser1","mockuser2"],"comparativeSummary":"This is a mock developer comparison for testing purposes.","skillComparison":{"JavaScript":[85.0,65.0],"Python":[45.0,90.0],"DevOps":[35.0,80.0]},"projectScale":{"mockuser1":["small","medium","medium"],"mockuser2":["medium","large","enterprise"]},"collaborationStyle":{"mockuser1":"Regular small contributions with detailed documentation","mockuser2":"Large feature implementations with minimal documentation"},"learningTrajectory":{"mockuser1":75.0,"mockuser2":85.0},"teamFitAnalysis":"These developers would complement each other well with their different skill sets.","complementarity":{"Frontend":["mockuser1"],"Backend":["mockuser2"],"Testing":["mockuser1","mockuser2"]}}`
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

// AnalyzeDeveloper analyzes a GitHub developer profile
func (s *LLMService) AnalyzeDeveloper(username string, repoData string, contributionData string, languageData string) (string, error) {
	// Format the prompt using the developer analysis template
	prompt := s.FormatPrompt(promptTemplates["developer_analysis"], map[string]string{
		"username":          username,
		"repo_data":         repoData,
		"contribution_data": contributionData,
		"language_data":     languageData,
	})

	// Call the LLM
	response, err := s.CallGeminiLLM(prompt)
	if err != nil {
		return "", fmt.Errorf("Developer Analysis Error: %v", err)
	}

	return response, nil
}

// CompareDevelopers compares multiple GitHub developer profiles
func (s *LLMService) CompareDevelopers(developerProfiles string, focus string) (string, error) {
	// Format the prompt using the developer comparison template
	prompt := s.FormatPrompt(promptTemplates["developer_comparison"], map[string]string{
		"developer_profiles": developerProfiles,
		"focus":              focus,
	})

	// Call the LLM
	response, err := s.CallGeminiLLM(prompt)
	if err != nil {
		return "", fmt.Errorf("Developer Comparison Error: %v", err)
	}

	return response, nil
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
Generate a Mermaid flowchart diagram that accurately represents the file hierarchy in a clean, organized manner (LR direction - left to right). The diagram should adhere to the following requirements:
1. **Graph Layout**: Use 'graph LR' for a left-to-right flow. You should branch out to the right for each level of the hierarchy. For each file under a directory, connect it to the directory node with a horizontal line. If two parent directories have a common child, that are not the common child, they are meant to be differnet files/folder, connect them to their respective parent directories. Meaning, duplicate file name are allowed, don't messed up the folder/files directory.
2. **Hierarchy**: Maintain proper parent-child relationships, ensuring each folder is connected only to its immediate parent, with no misplaced or disconnected nodes. Place level 0 folders at line horizontal line, level 1 folders or files at the next line, and so on. Make it easy for reader to identify the level of the folder/file.
3. **Node Shapes**:
- Use **rounded rectangles** for directories: 'folderName([folderName])'
- Use **rectangles** for files: 'fileName[fileName]'.
4. **Styling**:
- Apply visual distinction using these classes:
- Directories: 'classDef directory fill:#f9f,stroke:#333,stroke-width:1px;'
- Files: 'classDef file fill:#fff,stroke:#333,stroke-width:1px;'
- Ensure all directory nodes are styled with the 'directory' class and all file nodes with the 'file' class.
5. **Neatness**:
- Arrange nodes neatly with minimal crossing lines.
- Add suitable spacing between nodes to enhance readability.
6. **Output Format**:
- If files/folders have the same name, add a number to the end of the name (e.g., 'src', 'src1', 'src2', etc.).
- Do NOT include markdown code block syntax (three backticks).
- Return ONLY the raw Mermaid diagram code.
7. **Node Labels**:
- Display each file/folder using only its name (e.g., 'src', 'index.html', etc.), not the full path.
8. **Indentation**:
- Use proper indentation in the Mermaid code for better maintainability.
9. **Lines**:
- Arrange the lines neatly, use the way that can take up the least space.
`,
	
	/* New templates for developer analysis */
	"developer_analysis": `
Analyze the following GitHub user profile as a technical developer assessment:

Username: {username}

## Repository Data
{repo_data}

## Contribution Data
{contribution_data}

## Language Data
{language_data}

Provide a comprehensive developer profile analysis including:

1. Executive summary of the developer's technical capabilities (2-3 sentences)
2. Programming language proficiency analysis with estimated skill levels and years of experience
3. Project complexity assessment categorizing projects as toy/learning, small, medium, large, or enterprise-level
4. Contribution patterns and collaboration style assessment
5. Learning velocity and technology adoption timeline
6. Technical specialization areas and domain expertise identification
7. Technical skill assessment across: code quality, documentation, testing, performance, security, collaboration
8. Key technical strengths (top 3-5)
9. Areas for potential growth (top 3-5)

Return the analysis as a structured JSON with the following format exactly, with no extra text:
{
  "username": string,
  "executiveSummary": string,
  "languageAnalysis": [
    {
      "name": string,
      "percentage": float,
      "lines": int,
      "proficiency": string,
      "experience": float
    }
  ],
  "projectAnalysis": [
    {
      "name": string,
      "size": int,
      "stars": int,
      "forks": int,
      "complexityRank": string,
      "score": float
    }
  ],
  "contributionPatterns": {
    "totalCommits": int,
    "averageCommitsMonth": float,
    "consistencyScore": float,
    "prAcceptanceRate": float,
    "reviewParticipation": float,
    "monthlyActivity": [int]
  },
  "learningVelocity": {
    "adaptabilityScore": float,
    "technologyAdoption": {string: string},
    "growthTrajectory": string,
    "complexityTrend": [{string: float}]
  },
  "specializations": {
    "primaryDomains": [string],
    "secondaryDomains": [string],
    "frameworkExpertise": {string: float},
    "topicFocus": {string: float}
  },
  "skillAssessment": {
    "codeQuality": float,
    "documentation": float,
    "testing": float,
    "performance": float,
    "security": float,
    "collaboration": float
  },
  "technicalStrengths": [string],
  "growthAreas": [string]
}
`,
	"developer_comparison": `
Compare the following GitHub developer profiles:

{developer_profiles}

Provide a comprehensive comparison including:

1. Comparative summary of the developers' technical capabilities
2. Skill comparison across key technical areas
3. Project scale and complexity comparison
4. Collaboration style differences
5. Learning trajectory and adaptability comparison
6. Team fit analysis - how these developers would work together
7. Complementary skill sets analysis

Return the comparison as a structured JSON with the following format exactly, with no extra text:
{
  "developers": [string],
  "comparativeSummary": string,
  "skillComparison": {string: [float]},
  "projectScale": {string: [string]},
  "collaborationStyle": {string: string},
  "learningTrajectory": {string: float},
  "teamFitAnalysis": string,
  "complementarity": {string: [string]}
}
`,
}