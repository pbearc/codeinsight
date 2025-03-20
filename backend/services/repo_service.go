package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pbearc/codeinsight/models"
)

// RepoService provides methods for analyzing repositories
type RepoService struct {
	githubService *GitHubService
}

// NewRepoService creates a new repository service
func NewRepoService(githubService *GitHubService) *RepoService {
	return &RepoService{
		githubService: githubService,
	}
}

// parseGitHubURL parses a GitHub URL to extract owner and repo name
func parseGitHubURL(githubURL string) (string, string, error) {
	// Check if it's a valid URL
	parsedURL, err := url.Parse(githubURL)
	if err != nil {
		return "", "", err
	}
	
	// GitHub URL patterns
	githubRegexes := []*regexp.Regexp{
		// https://github.com/owner/repo
		regexp.MustCompile(`^/([^/]+)/([^/]+)/?$`),
		// https://github.com/owner/repo/blob/branch/path/to/file
		regexp.MustCompile(`^/([^/]+)/([^/]+)/blob/`),
		// https://github.com/owner/repo/tree/branch/path/to/directory
		regexp.MustCompile(`^/([^/]+)/([^/]+)/tree/`),
	}
	
	// Check if hostname is github.com
	if parsedURL.Hostname() != "github.com" {
		return "", "", fmt.Errorf("not a valid GitHub URL")
	}
	
	// Try to match the URL path with the patterns
	for _, regex := range githubRegexes {
		matches := regex.FindStringSubmatch(parsedURL.Path)
		if matches != nil && len(matches) >= 3 {
			return matches[1], matches[2], nil
		}
	}
	
	return "", "", fmt.Errorf("could not extract owner and repo from GitHub URL")
}

// AnalyzeGitHubRepository analyzes a GitHub repository structure with recursive traversal
func (s *RepoService) AnalyzeGitHubRepository(repoURL string) (*models.RepoStructure, error) {
	// Parse the GitHub URL to extract owner and repo
	owner, repo, err := parseGitHubURL(repoURL)
	if err != nil {
		return nil, fmt.Errorf("invalid GitHub URL: %w", err)
	}
	
	// Initialize repository structure
	repoStructure := &models.RepoStructure{
		Files:       []string{},
		Directories: []string{},
		FileSummary: make(map[string]string),
		Stats:       make(map[string]int),
	}
	
	// Initialize counters
	repoStructure.Stats["total_files"] = 0
	repoStructure.Stats["total_dirs"] = 0
	
	// Start recursive traversal from the root
	err = s.traverseGitHubDirectory(owner, repo, "", repoStructure)
	if err != nil {
		return nil, fmt.Errorf("failed to traverse repository: %w", err)
	}
	
	return repoStructure, nil
}

// traverseGitHubDirectory recursively traverses directories in a GitHub repository
func (s *RepoService) traverseGitHubDirectory(owner, repo, path string, repoStructure *models.RepoStructure) error {
	// Get contents of the current directory
	endpoint := fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, path)
	body, err := s.githubService.makeRequest("GET", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to get repository contents for path %s: %w", path, err)
	}
	
	// Parse response
	var contents []struct {
		Name        string `json:"name"`
		Path        string `json:"path"`
		Type        string `json:"type"`
		SHA         string `json:"sha"`
		Size        int    `json:"size"`
		URL         string `json:"url"`
		DownloadURL string `json:"download_url"`
	}
	
	if err := json.Unmarshal(body, &contents); err != nil {
		return fmt.Errorf("failed to parse repository contents: %w", err)
	}
	
	// Process files and directories
	for _, item := range contents {
		// Skip large binary files and unwanted directories
		if shouldSkipItem(item.Name, item.Size) {
			continue
		}
		
		if item.Type == "file" {
			// Process file
			repoStructure.Files = append(repoStructure.Files, item.Path)
			repoStructure.Stats["total_files"]++
			
			// Get file extension
			ext := filepath.Ext(item.Name)
			if ext != "" {
				ext = strings.ToLower(ext[1:]) // Remove the dot and make lowercase
				repoStructure.Stats["ext_"+ext] = repoStructure.Stats["ext_"+ext] + 1
			}
			
			// Get file content for important files
			if isImportantFile(item.Name) && item.Size < 500000 { // Skip files larger than 500KB
				content, err := s.githubService.GetFileContent(owner, repo, item.Path)
				if err == nil {
					// Generate a brief summary of the file
					summary := generateFileSummary(item.Name, content)
					repoStructure.FileSummary[item.Path] = summary
				}
			}
		} else if item.Type == "dir" {
			// Process directory
			repoStructure.Directories = append(repoStructure.Directories, item.Path)
			repoStructure.Stats["total_dirs"]++
			
			// Recursively process the subdirectory
			err = s.traverseGitHubDirectory(owner, repo, item.Path, repoStructure)
			if err != nil {
				// Log the error but continue with other directories
				fmt.Printf("Error traversing directory %s: %v\n", item.Path, err)
			}
		}
	}
	
	return nil
}

// shouldSkipItem determines if a file or directory should be skipped
func shouldSkipItem(name string, size int) bool {
	// Skip dot directories
	if name == ".git" || name == ".github" || name == "node_modules" {
		return true
	}
	
	// Skip very large files
	if size > 10000000 { // 10MB
		return true
	}
	
	// Skip binary/media files
	binaryExtensions := []string{
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".ico", ".svg",
		".mp3", ".mp4", ".avi", ".mov", ".wmv",
		".zip", ".tar", ".gz", ".7z", ".rar",
		".exe", ".dll", ".so", ".dylib",
		".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
	}
	
	ext := strings.ToLower(filepath.Ext(name))
	for _, binExt := range binaryExtensions {
		if ext == binExt {
			return true
		}
	}
	
	return false
}

// AnalyzeLocalDirectory analyzes a local directory structure
func (s *RepoService) AnalyzeLocalDirectory(dirPath string) (*models.RepoStructure, error) {
	// Check if directory exists
	info, err := os.Stat(dirPath)
	if err != nil {
		return nil, fmt.Errorf("directory error: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", dirPath)
	}
	
	// Initialize repository structure
	repoStructure := &models.RepoStructure{
		Files:       []string{},
		Directories: []string{},
		FileSummary: make(map[string]string),
		Stats:       make(map[string]int),
	}
	
	// Count file types
	repoStructure.Stats["total_files"] = 0
	repoStructure.Stats["total_dirs"] = 0
	
	// Walk through the directory
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Get relative path
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}
		
		// Skip .git and node_modules directories
		if (info.IsDir() && (info.Name() == ".git" || info.Name() == "node_modules")) {
			return filepath.SkipDir
		}
		
		// Process files and directories
		if info.IsDir() {
			if relPath != "." { // Skip the root directory
				repoStructure.Directories = append(repoStructure.Directories, relPath)
				repoStructure.Stats["total_dirs"]++
			}
		} else {
			repoStructure.Files = append(repoStructure.Files, relPath)
			repoStructure.Stats["total_files"]++
			
			// Get file extension
			ext := filepath.Ext(info.Name())
			if ext != "" {
				ext = strings.ToLower(ext[1:]) // Remove the dot and make lowercase
				repoStructure.Stats["ext_"+ext] = repoStructure.Stats["ext_"+ext] + 1
			}
			
			// Get file content for important files
			if isImportantFile(info.Name()) && info.Size() < 500000 { // Skip files larger than 500KB
				content, err := ioutil.ReadFile(path)
				if err == nil {
					// Generate a brief summary of the file
					summary := generateFileSummary(info.Name(), string(content))
					repoStructure.FileSummary[relPath] = summary
				}
			}
		}
		
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}
	
	return repoStructure, nil
}

// FormatRepositoryStructure formats repository structure for LLM processing
func (s *RepoService) FormatRepositoryStructure(repoStructure *models.RepoStructure) string {
	var sb strings.Builder
	
	// Write repository statistics
	sb.WriteString("# Repository Structure\n\n")
	sb.WriteString("## Statistics\n")
	sb.WriteString(fmt.Sprintf("- Total files: %d\n", repoStructure.Stats["total_files"]))
	sb.WriteString(fmt.Sprintf("- Total directories: %d\n", repoStructure.Stats["total_dirs"]))
	
	// Write file types
	sb.WriteString("\n## File Types\n")
	for key, value := range repoStructure.Stats {
		if strings.HasPrefix(key, "ext_") {
			extension := strings.TrimPrefix(key, "ext_")
			sb.WriteString(fmt.Sprintf("- %s: %d files\n", extension, value))
		}
	}
	
	// Write directory structure
	sb.WriteString("\n## Directory Structure\n")
	for _, dir := range repoStructure.Directories {
		sb.WriteString(fmt.Sprintf("- %s/\n", dir))
	}
	
	// Write file structure
	sb.WriteString("\n## Files\n")
	for _, file := range repoStructure.Files {
		sb.WriteString(fmt.Sprintf("- %s\n", file))
	}
	
	// Write file summaries
	sb.WriteString("\n## File Summaries\n")
	for path, summary := range repoStructure.FileSummary {
		sb.WriteString(fmt.Sprintf("### %s\n", path))
		sb.WriteString(fmt.Sprintf("%s\n\n", summary))
	}
	
	return sb.String()
}

// isImportantFile checks if a file is important for analysis
func isImportantFile(filename string) bool {
	// List of important file patterns
	importantPatterns := []string{
		// Configuration files
		"package.json", "composer.json", "Gemfile", "requirements.txt",
		"setup.py", "pom.xml", "build.gradle", "Cargo.toml", "go.mod",
		".gitignore", ".dockerignore", "Dockerfile", "docker-compose.yml",
		
		// Documentation
		"README", "LICENSE", "CONTRIBUTING", "CHANGELOG",
		
		// Source code
		".js", ".ts", ".jsx", ".tsx", ".py", ".java", ".go", ".rb", ".php",
		".c", ".cpp", ".h", ".hpp", ".cs", ".scala", ".swift", ".kt",
	}
	
	// Check if filename matches any important pattern
	for _, pattern := range importantPatterns {
		if strings.HasSuffix(strings.ToUpper(filename), strings.ToUpper(pattern)) {
			return true
		}
	}
	
	return false
}

// generateFileSummary generates a brief summary of a file's content
func generateFileSummary(filename string, content string) string {
	// Get file extension
	ext := filepath.Ext(filename)
	
	// Determine file type
	var fileType string
	switch strings.ToLower(ext) {
	case ".js", ".jsx", ".ts", ".tsx":
		fileType = "JavaScript/TypeScript"
	case ".py":
		fileType = "Python"
	case ".java":
		fileType = "Java"
	case ".go":
		fileType = "Go"
	case ".rb":
		fileType = "Ruby"
	case ".php":
		fileType = "PHP"
	case ".json":
		fileType = "JSON"
	case ".md":
		fileType = "Markdown"
	default:
		fileType = "Text"
	}
	
	// Truncate content if too long
	maxLength := 1000
	if len(content) > maxLength {
		content = content[:maxLength] + "..."
	}
	
	// For code files, try to extract important parts
	var summary string
	switch fileType {
	case "JavaScript/TypeScript":
		summary = extractJavaScriptSummary(content)
	case "Python":
		summary = extractPythonSummary(content)
	case "Java":
		summary = extractJavaSummary(content)
	case "Go":
		summary = extractGoSummary(content)
	default:
		// For other file types, just take the first few lines
		lines := strings.Split(content, "\n")
		if len(lines) > 10 {
			summary = strings.Join(lines[:10], "\n") + "\n..."
		} else {
			summary = content
		}
	}
	
	return fmt.Sprintf("**File Type**: %s\n\n**Content**:\n```\n%s\n```", fileType, summary)
}

// extractJavaScriptSummary extracts a summary from JavaScript/TypeScript code
func extractJavaScriptSummary(content string) string {
	// Extract imports/requires
	importRegex := regexp.MustCompile(`(import .+?from .+?;|const .+? = require\(.+?\))`)
	imports := importRegex.FindAllString(content, -1)
	
	// Extract function and class declarations
	funcRegex := regexp.MustCompile(`(function\s+\w+\s*\([^)]*\)|const\s+\w+\s*=\s*\([^)]*\)\s*=>|class\s+\w+)`)
	functions := funcRegex.FindAllString(content, -1)
	
	// Build summary
	var summary strings.Builder
	
	// Add imports
	if len(imports) > 0 {
		summary.WriteString("// Imports\n")
		for i, imp := range imports {
			if i < 5 { // Limit to 5 imports
				summary.WriteString(imp + "\n")
			} else {
				summary.WriteString("// ... more imports\n")
				break
			}
		}
		summary.WriteString("\n")
	}
	
	// Add functions/classes
	if len(functions) > 0 {
		summary.WriteString("// Functions/Classes\n")
		for i, fn := range functions {
			if i < 5 { // Limit to 5 functions
				summary.WriteString(fn + " { ... }\n")
			} else {
				summary.WriteString("// ... more functions/classes\n")
				break
			}
		}
	}
	
	return summary.String()
}

// extractPythonSummary extracts a summary from Python code
func extractPythonSummary(content string) string {
	// Extract imports
	importRegex := regexp.MustCompile(`(from .+? import .+|import .+)`)
	imports := importRegex.FindAllString(content, -1)
	
	// Extract function and class declarations
	funcRegex := regexp.MustCompile(`(def\s+\w+\s*\([^)]*\)|class\s+\w+)`)
	functions := funcRegex.FindAllString(content, -1)
	
	// Build summary
	var summary strings.Builder
	
	// Add imports
	if len(imports) > 0 {
		summary.WriteString("# Imports\n")
		for i, imp := range imports {
			if i < 5 { // Limit to 5 imports
				summary.WriteString(imp + "\n")
			} else {
				summary.WriteString("# ... more imports\n")
				break
			}
		}
		summary.WriteString("\n")
	}
	
	// Add functions/classes
	if len(functions) > 0 {
		summary.WriteString("# Functions/Classes\n")
		for i, fn := range functions {
			if i < 5 { // Limit to 5 functions
				summary.WriteString(fn + ":\n    ...\n")
			} else {
				summary.WriteString("# ... more functions/classes\n")
				break
			}
		}
	}
	
	return summary.String()
}

// extractJavaSummary extracts a summary from Java code
func extractJavaSummary(content string) string {
	// Extract imports
	importRegex := regexp.MustCompile(`import .+;`)
	imports := importRegex.FindAllString(content, -1)
	
	// Extract class and method declarations
	classRegex := regexp.MustCompile(`(public|private|protected)?\s*(class|interface|enum)\s+\w+`)
	classes := classRegex.FindAllString(content, -1)
	
	methodRegex := regexp.MustCompile(`(public|private|protected)?\s*(static)?\s*\w+\s+\w+\s*\([^)]*\)`)
	methods := methodRegex.FindAllString(content, -1)
	
	// Build summary
	var summary strings.Builder
	
	// Add imports
	if len(imports) > 0 {
		summary.WriteString("// Imports\n")
		for i, imp := range imports {
			if i < 5 { // Limit to 5 imports
				summary.WriteString(imp + "\n")
			} else {
				summary.WriteString("// ... more imports\n")
				break
			}
		}
		summary.WriteString("\n")
	}
	
	// Add classes
	if len(classes) > 0 {
		summary.WriteString("// Classes\n")
		for i, cls := range classes {
			if i < 3 { // Limit to 3 classes
				summary.WriteString(cls + " { ... }\n")
			} else {
				summary.WriteString("// ... more classes\n")
				break
			}
		}
		summary.WriteString("\n")
	}
	
	// Add methods
	if len(methods) > 0 {
		summary.WriteString("// Methods\n")
		for i, mtd := range methods {
			if i < 5 { // Limit to 5 methods
				summary.WriteString(mtd + " { ... }\n")
			} else {
				summary.WriteString("// ... more methods\n")
				break
			}
		}
	}
	
	return summary.String()
}

// extractGoSummary extracts a summary from Go code
func extractGoSummary(content string) string {
	// Extract package and imports
	packageRegex := regexp.MustCompile(`package\s+\w+`)
	packageMatch := packageRegex.FindString(content)
	
	importRegex := regexp.MustCompile(`import\s+\(\s+(.+?)\s+\)`)
	importMatches := importRegex.FindAllStringSubmatch(content, -1)
	
	// Extract function declarations
	funcRegex := regexp.MustCompile(`func\s+(\w+|\([^)]+\)\s*\w+)\s*\([^)]*\)`)
	functions := funcRegex.FindAllString(content, -1)
	
	// Extract struct declarations
	structRegex := regexp.MustCompile(`type\s+\w+\s+struct\s*\{`)
	structs := structRegex.FindAllString(content, -1)
	
	// Build summary
	var summary strings.Builder
	
	// Add package
	if packageMatch != "" {
		summary.WriteString(packageMatch + "\n\n")
	}
	
	// Add imports
	if len(importMatches) > 0 {
		summary.WriteString("import (\n")
		for _, match := range importMatches {
			if len(match) > 1 {
				importLines := strings.Split(match[1], "\n")
				for i, line := range importLines {
					if i < 5 { // Limit to 5 imports
						summary.WriteString("\t" + strings.TrimSpace(line) + "\n")
					} else {
						summary.WriteString("\t// ... more imports\n")
						break
					}
				}
			}
		}
		summary.WriteString(")\n\n")
	}
	
	// Add structs
	if len(structs) > 0 {
		summary.WriteString("// Structs\n")
		for i, str := range structs {
			if i < 3 { // Limit to 3 structs
				summary.WriteString(str + " ... }\n")
			} else {
				summary.WriteString("// ... more structs\n")
				break
			}
		}
		summary.WriteString("\n")
	}
	
	// Add functions
	if len(functions) > 0 {
		summary.WriteString("// Functions\n")
		for i, fn := range functions {
			if i < 5 { // Limit to 5 functions
				summary.WriteString(fn + " { ... }\n")
			} else {
				summary.WriteString("// ... more functions\n")
				break
			}
		}
	}
	
	return summary.String()
}