package models

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
	Inline      bool   `json:"inline" default:"false"`
}

// ReadmeGeneratorRequest represents a request to generate a README file
type ReadmeGeneratorRequest struct {
	RepoURL     string `json:"repoUrl"`
	FolderPath  string `json:"folderPath"`
	ProjectName string `json:"projectName"`
	Description string `json:"description"`
}

// ProjectVisualizationRequest represents a request to visualize a project structure
type ProjectVisualizationRequest struct {
	RepoURL    string `json:"repoUrl"`
	FolderPath string `json:"folderPath"`
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
	Inline        bool   `json:"inline"`
}

// ReadmeResponse represents a generated README file
type ReadmeResponse struct {
	Content string `json:"content"`
}

// ProjectVisualizationResponse represents a project visualization diagram
type ProjectVisualizationResponse struct {
	DiagramCode string `json:"diagramCode"`
	Type        string `json:"type"`
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

// RepoStructure represents a repository structure for scanning
type RepoStructure struct {
	Files       []string            `json:"files"`
	Directories []string            `json:"directories"`
	FileSummary map[string]string   `json:"fileSummary"`
	Stats       map[string]int      `json:"stats"`
}

// DeveloperAnalysisRequest represents a request to analyze a GitHub developer profile
type DeveloperAnalysisRequest struct {
	Username string `json:"username" binding:"required"`
	Depth    string `json:"depth" default:"medium"` // Options: light, medium, full
}

// DeveloperCompareRequest represents a request to compare multiple GitHub developer profiles
type DeveloperCompareRequest struct {
	Usernames []string `json:"usernames" binding:"required,min=2"`
	Focus     string   `json:"focus" default:"general"` // Options: general, languages, projects, contributions
}

// LanguageProficiency represents a programming language and associated metrics
type LanguageProficiency struct {
	Name        string  `json:"name"`
	Percentage  float64 `json:"percentage"`
	Lines       int     `json:"lines"`
	Proficiency string  `json:"proficiency"` // beginner, intermediate, advanced, expert
	Experience  float64 `json:"experience"`  // estimated years
}

// ProjectComplexity represents the complexity analysis of a project
type ProjectComplexity struct {
	Name           string  `json:"name"`
	Size           int     `json:"size"`
	Stars          int     `json:"stars"`
	Forks          int     `json:"forks"`
	ComplexityRank string  `json:"complexityRank"` // toy, small, medium, large, enterprise
	Score          float64 `json:"score"`
}

// ContributionPattern represents analysis of GitHub contribution patterns
type ContributionPattern struct {
	TotalCommits        int     `json:"totalCommits"`
	AverageCommitsMonth float64 `json:"averageCommitsMonth"`
	ConsistencyScore    float64 `json:"consistencyScore"` // 0-100
	PRAcceptanceRate    float64 `json:"prAcceptanceRate"` // percentage
	ReviewParticipation float64 `json:"reviewParticipation"` // 0-100
	MonthlyActivity     []int   `json:"monthlyActivity"` // Last 12 months
}

// LearningVelocity represents analysis of technology adoption and growth
type LearningVelocity struct {
	AdaptabilityScore float64                `json:"adaptabilityScore"` // 0-100
	TechnologyAdoption map[string]string     `json:"technologyAdoption"` // tech name -> approx time
	GrowthTrajectory  string                 `json:"growthTrajectory"` // steady, accelerating, plateaued
	ComplexityTrend   []map[string]float64   `json:"complexityTrend"` // timeline of complexity metrics
}

// Specialization represents domain expertise identification
type Specialization struct {
	PrimaryDomains     []string            `json:"primaryDomains"` // e.g., "web development", "data science"
	SecondaryDomains   []string            `json:"secondaryDomains"`
	FrameworkExpertise map[string]float64  `json:"frameworkExpertise"` // framework -> proficiency (0-100)
	TopicFocus         map[string]float64  `json:"topicFocus"` // topic -> focus percentage
}

// SkillAssessment represents a radar chart of technical skills
type SkillAssessment struct {
	CodeQuality    float64 `json:"codeQuality"` // 0-100
	Documentation  float64 `json:"documentation"` // 0-100
	Testing        float64 `json:"testing"` // 0-100
	Performance    float64 `json:"performance"` // 0-100
	Security       float64 `json:"security"` // 0-100
	Collaboration  float64 `json:"collaboration"` // 0-100
}

// DeveloperProfile represents a complete GitHub developer profile analysis
type DeveloperProfile struct {
	Username            string               `json:"username"`
	ExecutiveSummary    string               `json:"executiveSummary"`
	LanguageAnalysis    []LanguageProficiency `json:"languageAnalysis"`
	ProjectAnalysis     []ProjectComplexity   `json:"projectAnalysis"`
	ContributionPatterns ContributionPattern  `json:"contributionPatterns"`
	LearningVelocity    LearningVelocity     `json:"learningVelocity"`
	Specializations     Specialization       `json:"specializations"`
	SkillAssessment     SkillAssessment      `json:"skillAssessment"`
	TechnicalStrengths  []string             `json:"technicalStrengths"`
	GrowthAreas         []string             `json:"growthAreas"`
}

// ComparativeAnalysis represents a comparison between multiple developers
type ComparativeAnalysis struct {
	Developers         []string                `json:"developers"`
	ComparativeSummary string                  `json:"comparativeSummary"`
	SkillComparison    map[string][]float64    `json:"skillComparison"` // skill -> scores array
	ProjectScale       map[string][]string     `json:"projectScale"` // developer -> project scale distribution
	CollaborationStyle map[string]string       `json:"collaborationStyle"` // developer -> style description
	LearningTrajectory map[string]float64      `json:"learningTrajectory"` // developer -> adaptability score
	TeamFitAnalysis    string                  `json:"teamFitAnalysis"`
	Complementarity    map[string][]string     `json:"complementarity"` // skill -> developers with strength
}

// DeveloperAnalysisResponse represents a developer analysis result from the LLM
type DeveloperAnalysisResponse struct {
	Success bool             `json:"success"`
	Profile DeveloperProfile `json:"profile"`
}

// DeveloperCompareResponse represents a developer comparison result from the LLM
type DeveloperCompareResponse struct {
	Success   bool                `json:"success"`
	Analysis  ComparativeAnalysis `json:"analysis"`
}