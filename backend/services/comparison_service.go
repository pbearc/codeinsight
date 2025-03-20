package services

import (
	"fmt"
	"regexp"
	"strings"
)

// ComparisonService provides methods for comparing code implementations
type ComparisonService struct{}

// NewComparisonService creates a new comparison service
func NewComparisonService() *ComparisonService {
	return &ComparisonService{}
}

// CompareImplementations compares multiple implementations of a function/algorithm
func (s *ComparisonService) CompareImplementations(implementations []map[string]interface{}, language string) (string, error) {
	if len(implementations) < 2 {
		return "", fmt.Errorf("at least 2 implementations required for comparison")
	}
	
	// Extract code from implementations
	codes := make([]string, len(implementations))
	sources := make([]string, len(implementations))
	
	for i, impl := range implementations {
		file, ok := impl["file"].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("invalid implementation format at index %d", i)
		}
		
		content, ok := file["content"].(string)
		if !ok {
			return "", fmt.Errorf("missing code content at index %d", i)
		}
		
		repo, ok := impl["repository"].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("invalid repository data at index %d", i)
		}
		
		fullName, _ := repo["full_name"].(string)
		
		codes[i] = content
		sources[i] = fullName
	}
	
	// Generate comparison based on language
	var comparison string
	
	switch strings.ToLower(language) {
	case "python":
		comparison = s.comparePythonImplementations(codes, sources)
	case "javascript":
		comparison = s.compareJavaScriptImplementations(codes, sources)
	case "java":
		comparison = s.compareJavaImplementations(codes, sources)
	default:
		comparison = s.compareGenericImplementations(codes, sources, language)
	}
	
	return comparison, nil
}

// comparePythonImplementations generates a comparison of Python implementations
func (s *ComparisonService) comparePythonImplementations(codes []string, sources []string) string {
	markdown := "## Implementation Comparison\n\n"
	
	// Add overview
	markdown += "### Overview\n\n"
	markdown += fmt.Sprintf("Comparing %d different Python implementations:\n\n", len(codes))
	
	for i, source := range sources {
		markdown += fmt.Sprintf("%d. From repository: `%s`\n", i+1, source)
	}
	
	// Analyze implementations
	var (
		recursiveCount  int
		inPlaceCount    int
		docstringCount  int
		optimizedCount  int
		listCompCount   int
		randomPivotCount int
	)
	
	for _, code := range codes {
		// Check if recursive
		if strings.Contains(code, "return quicksort(") {
			recursiveCount++
		}
		
		// Check if in-place
		if strings.Contains(code, "def partition") {
			inPlaceCount++
		}
		
		// Check if well-documented
		if strings.Contains(code, "\"\"\"") || strings.Contains(code, "Time Complexity") {
			docstringCount++
		}
		
		// Check if using optimization
		if strings.Contains(code, "insertion_sort") || 
		   regexp.MustCompile(`if\s+len\(.*\)\s*<=\s*\d+`).MatchString(code) {
			optimizedCount++
		}
		
		// Check if using list comprehension
		if strings.Contains(code, "[x for x in") {
			listCompCount++
		}
		
		// Check if using random pivot
		if strings.Contains(code, "random") {
			randomPivotCount++
		}
	}
	
	// Approach comparison
	markdown += "\n### Implementation Approaches\n\n"
	
	totalImpls := len(codes)
	markdown += fmt.Sprintf("- **Recursive implementations**: %d/%d\n", recursiveCount, totalImpls)
	markdown += fmt.Sprintf("- **In-place implementations**: %d/%d\n", inPlaceCount, totalImpls)
	markdown += fmt.Sprintf("- **Well-documented code**: %d/%d\n", docstringCount, totalImpls)
	markdown += fmt.Sprintf("- **With optimizations**: %d/%d\n", optimizedCount, totalImpls)
	markdown += fmt.Sprintf("- **Using list comprehensions**: %d/%d\n", listCompCount, totalImpls)
	markdown += fmt.Sprintf("- **Using randomized pivots**: %d/%d\n", randomPivotCount, totalImpls)
	
	// Recommendations
	markdown += "\n### Recommendations\n\n"
	
	if randomPivotCount > 0 {
		markdown += "- **For reliability**: Implementations with randomized pivots are more robust against worst-case scenarios\n"
	}
	
	if optimizedCount > 0 {
		markdown += "- **For performance**: Consider implementations with optimizations for small arrays\n"
	}
	
	if inPlaceCount > 0 && listCompCount > 0 {
		markdown += "- **Memory vs. readability tradeoff**: In-place implementations are more memory-efficient but list comprehension versions are more readable\n"
	}
	
	// Time & space complexity
	markdown += "\n### Complexity Analysis\n\n"
	markdown += "- **Time Complexity**: O(n log n) average case, O(n²) worst case\n"
	
	if inPlaceCount > 0 {
		markdown += "- **Space Complexity**: O(log n) for in-place implementations, O(n) for others\n"
	} else {
		markdown += "- **Space Complexity**: O(n) for all implementations\n"
	}
	
	return markdown
}

// compareJavaScriptImplementations generates a comparison of JavaScript implementations
func (s *ComparisonService) compareJavaScriptImplementations(codes []string, sources []string) string {
	markdown := "## Implementation Comparison\n\n"
	
	// Add overview
	markdown += "### Overview\n\n"
	markdown += fmt.Sprintf("Comparing %d different JavaScript implementations:\n\n", len(codes))
	
	for i, source := range sources {
		markdown += fmt.Sprintf("%d. From repository: `%s`\n", i+1, source)
	}
	
	// Analyze implementations
	var (
		es6Count        int
		inPlaceCount    int
		commentedCount  int
		optimizedCount  int
		recursiveCount  int
	)
	
	for _, code := range codes {
		// Check if using ES6 features
		if strings.Contains(code, "=>") || strings.Contains(code, "...") || 
		   strings.Contains(code, "const") {
			es6Count++
		}
		
		// Check if in-place
		if strings.Contains(code, "function partition") || 
		   regexp.MustCompile(`\w+\s*\(\s*arr\s*,\s*left\s*,\s*right`).MatchString(code) {
			inPlaceCount++
		}
		
		// Check if well-commented
		if strings.Contains(code, "/**") || strings.Contains(code, "Time Complexity") ||
		   strings.Contains(code, "@param") {
			commentedCount++
		}
		
		// Check if using optimization
		if strings.Contains(code, "insertionSort") || 
		   regexp.MustCompile(`if\s*\(\s*arr\.length\s*<=\s*\d+`).MatchString(code) {
			optimizedCount++
		}
		
		// Check if recursive
		if strings.Contains(code, "return quicksort") {
			recursiveCount++
		}
	}
	
	// Approach comparison
	markdown += "\n### Implementation Approaches\n\n"
	
	totalImpls := len(codes)
	markdown += fmt.Sprintf("- **Using modern ES6+ features**: %d/%d\n", es6Count, totalImpls)
	markdown += fmt.Sprintf("- **In-place implementations**: %d/%d\n", inPlaceCount, totalImpls)
	markdown += fmt.Sprintf("- **Well-documented code**: %d/%d\n", commentedCount, totalImpls)
	markdown += fmt.Sprintf("- **With optimizations**: %d/%d\n", optimizedCount, totalImpls)
	markdown += fmt.Sprintf("- **Recursive implementations**: %d/%d\n", recursiveCount, totalImpls)
	
	// Recommendations
	markdown += "\n### Recommendations\n\n"
	
	if es6Count > 0 {
		markdown += "- **For modern codebases**: Prefer the ES6+ implementations for better readability\n"
	}
	
	if optimizedCount > 0 {
		markdown += "- **For performance**: Consider implementations with optimizations for small arrays\n"
	}
	
	if inPlaceCount > 0 {
		markdown += "- **For memory efficiency**: In-place implementations use less memory\n"
	}
	
	// Time & space complexity
	markdown += "\n### Complexity Analysis\n\n"
	markdown += "- **Time Complexity**: O(n log n) average case, O(n²) worst case\n"
	
	if inPlaceCount > 0 {
		markdown += "- **Space Complexity**: O(log n) for in-place implementations, O(n) for others\n"
	} else {
		markdown += "- **Space Complexity**: O(n) for all implementations\n"
	}
	
	return markdown
}

// compareJavaImplementations generates a comparison of Java implementations
func (s *ComparisonService) compareJavaImplementations(codes []string, sources []string) string {
	markdown := "## Implementation Comparison\n\n"
	
	// Add overview
	markdown += "### Overview\n\n"
	markdown += fmt.Sprintf("Comparing %d different Java implementations:\n\n", len(codes))
	
	for i, source := range sources {
		markdown += fmt.Sprintf("%d. From repository: `%s`\n", i+1, source)
	}
	
	// Analyze implementations
	var (
		inPlaceCount    int
		genericCount    int
		commentedCount  int
		optimizedCount  int
	)
	
	for _, code := range codes {
		// Check if in-place
		if strings.Contains(code, "private static int partition") || 
		   strings.Contains(code, "swap(") {
			inPlaceCount++
		}
		
		// Check if generic
		if strings.Contains(code, "<T") || strings.Contains(code, "Comparable") {
			genericCount++
		}
		
		// Check if well-commented
		if strings.Contains(code, "/**") || strings.Contains(code, "@param") ||
		   strings.Contains(code, "Time complexity") {
			commentedCount++
		}
		
		// Check if using optimization
		if strings.Contains(code, "insertionSort") || 
		   regexp.MustCompile(`if\s*\(\s*.*\.length\s*<=\s*\d+`).MatchString(code) {
			optimizedCount++
		}
	}
	
	// Approach comparison
	markdown += "\n### Implementation Approaches\n\n"
	
	totalImpls := len(codes)
	markdown += fmt.Sprintf("- **In-place implementations**: %d/%d\n", inPlaceCount, totalImpls)
	markdown += fmt.Sprintf("- **Generic implementations**: %d/%d\n", genericCount, totalImpls)
	markdown += fmt.Sprintf("- **Well-documented code**: %d/%d\n", commentedCount, totalImpls)
	markdown += fmt.Sprintf("- **With optimizations**: %d/%d\n", optimizedCount, totalImpls)
	
	// Recommendations
	markdown += "\n### Recommendations\n\n"
	
	if genericCount > 0 {
		markdown += "- **For reusability**: Generic implementations can work with any comparable type\n"
	}
	
	if optimizedCount > 0 {
		markdown += "- **For performance**: Consider implementations with optimizations for small arrays\n"
	}
	
	if inPlaceCount > 0 {
		markdown += "- **For memory efficiency**: In-place implementations use less memory\n"
	}
	
	// Time & space complexity
	markdown += "\n### Complexity Analysis\n\n"
	markdown += "- **Time Complexity**: O(n log n) average case, O(n²) worst case\n"
	markdown += "- **Space Complexity**: O(log n) for recursive calls in most implementations\n"
	
	return markdown
}

// compareGenericImplementations generates a generic comparison for other languages
func (s *ComparisonService) compareGenericImplementations(codes []string, sources []string, language string) string {
	markdown := "## Implementation Comparison\n\n"
	
	// Add overview
	markdown += "### Overview\n\n"
	markdown += fmt.Sprintf("Comparing %d different %s implementations:\n\n", len(codes), language)
	
	for i, source := range sources {
		markdown += fmt.Sprintf("%d. From repository: `%s`\n", i+1, source)
	}
	
	// General analysis
	markdown += "\n### General Analysis\n\n"
	markdown += "Key differences between implementations:\n\n"
	markdown += "- **Pivot selection strategy**: Different implementations may choose first, last, middle, or random elements as pivots\n"
	markdown += "- **In-place vs. auxiliary storage**: Some implementations modify the array in place, while others create new arrays\n"
	markdown += "- **Handling of duplicates**: Some implementations have special handling for duplicate elements\n"
	markdown += "- **Recursive vs. iterative**: Most implementations use recursion, but some may use an iterative approach with a stack\n"
	markdown += "- **Optimizations**: Some implementations switch to insertion sort for small arrays\n"
	
	// Code size comparison
	markdown += "\n### Code Size Comparison\n\n"
	markdown += "| Implementation | Lines of Code |\n"
	markdown += "| -------------- | ------------ |\n"
	
	for i, code := range codes {
		lineCount := len(strings.Split(code, "\n"))
		markdown += fmt.Sprintf("| Implementation %d | %d |\n", i+1, lineCount)
	}
	
	// Recommendations
	markdown += "\n### Recommendations\n\n"
	markdown += "- When choosing an implementation, consider the tradeoff between readability and efficiency\n"
	markdown += "- For production code, prefer well-tested implementations with good documentation\n"
	markdown += "- Consider the input characteristics (size, partially sorted, duplicates) when selecting an algorithm variant\n"
	
	// Time & space complexity
	markdown += "\n### Complexity Analysis\n\n"
	markdown += "- **Time Complexity**: O(n log n) average case, O(n²) worst case\n"
	markdown += "- **Space Complexity**: Varies between implementations (O(log n) to O(n))\n"
	
	return markdown
}