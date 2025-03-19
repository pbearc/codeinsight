// frontend/src/pages/LibraryExplorer.js
import React, { useState } from "react";
import axios from "axios";
import {
  Container,
  Row,
  Col,
  Form,
  Button,
  Alert,
  Card,
} from "react-bootstrap";
import ReactMarkdown from "react-markdown";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { vscDarkPlus } from "react-syntax-highlighter/dist/esm/styles/prism";

// Components
import Spinner from "../components/Spinner";
import LibraryUsageExample from "../components/LibraryUsageExample";

function LibraryExplorer() {
  const [library, setLibrary] = useState("");
  const [language, setLanguage] = useState("javascript");
  const [limit, setLimit] = useState(5);
  const [examples, setExamples] = useState([]);
  const [summary, setSummary] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const API_BASE_URL =
    process.env.REACT_APP_API_BASE_URL || "http://localhost:5000/api";

  const languageOptions = [
    { value: "javascript", label: "JavaScript" },
    { value: "python", label: "Python" },
    { value: "java", label: "Java" },
    { value: "csharp", label: "C#" },
    { value: "cpp", label: "C++" },
    { value: "php", label: "PHP" },
    { value: "ruby", label: "Ruby" },
    { value: "go", label: "Go" },
    { value: "typescript", label: "TypeScript" },
  ];

  // Handle form submission
  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!library.trim()) {
      setError("Please enter a library name");
      return;
    }

    try {
      setLoading(true);
      setError(null);
      setExamples([]);
      setSummary(null);

      // Fetch library usage examples
      const response = await axios.get(`${API_BASE_URL}/github/library-usage`, {
        params: {
          library,
          language,
          limit,
        },
      });

      if (response.data && response.data.success) {
        setExamples(response.data.data);

        // If examples were found, generate a summary using the LLM
        if (response.data.data.length > 0) {
          await generateSummary(response.data.data);
        } else {
          setError("No usage examples found for this library");
        }
      } else {
        throw new Error("Failed to fetch library usage examples");
      }
    } catch (err) {
      setError(err.message || "An error occurred");
    } finally {
      setLoading(false);
    }
  };

  // Generate a summary of the library usage patterns
  const generateSummary = async (examples) => {
    try {
      // Extract code snippets from examples
      const snippets = examples.map((ex) => ({
        repo: ex.repository.full_name,
        code: ex.file.content.slice(0, 1000), // Limit size for LLM
      }));

      // Generate summary using LLM
      const response = await axios.post(`${API_BASE_URL}/analyze`, {
        code: JSON.stringify(snippets),
        language,
        analysisType: "customPrompt",
        prompt: `
          Analyze these ${snippets.length} examples of using the "${library}" library in ${language}.
          Identify common usage patterns, best practices, and provide a summary of how this library is typically used.
          Include specific insights about:
          1. Common import/require patterns
          2. Most frequently used functions or methods
          3. Typical initialization patterns
          4. Error handling approaches
          5. Any other notable patterns or practices
          
          Format your response as markdown.
        `,
      });

      if (response.data && response.data.success) {
        setSummary(response.data.data.analysis);
      }
    } catch (err) {
      console.error("Failed to generate summary:", err);
      // We don't set an error here as we still want to show the examples
      // even if the summary generation fails
    }
  };

  // Popular library suggestions for each language
  const popularLibraries = {
    javascript: ["react", "express", "lodash", "axios", "moment"],
    python: ["numpy", "pandas", "requests", "django", "flask"],
    java: ["spring", "jackson", "gson", "guava", "apache commons"],
    typescript: ["angular", "nest", "typeorm", "rxjs", "zod"],
    csharp: [
      "newtonsoft.json",
      "entity framework",
      "automapper",
      "nlog",
      "moq",
    ],
    // Add more languages and libraries as needed
  };

  // Get suggestions for the current language
  const getCurrentSuggestions = () => {
    return popularLibraries[language] || popularLibraries.javascript;
  };

  // Handle clicking a suggestion
  const handleSuggestionClick = (suggestion) => {
    setLibrary(suggestion);
  };

  return (
    <Container className="py-4">
      <div className="text-center mb-4">
        <h1>Library Explorer</h1>
        <p className="lead">
          Discover how developers use libraries and frameworks across GitHub.
          Find common patterns and best practices.
        </p>
      </div>

      <div className="bg-white p-4 border rounded shadow-sm mb-4">
        <Form onSubmit={handleSubmit}>
          <Form.Group className="mb-3">
            <Form.Label>Library Name</Form.Label>
            <Form.Control
              type="text"
              placeholder="Enter a library or framework name (e.g., react, pandas)"
              value={library}
              onChange={(e) => setLibrary(e.target.value)}
              required
            />
          </Form.Group>

          <Row className="mb-3">
            <Col md={6}>
              <Form.Group>
                <Form.Label>Language</Form.Label>
                <Form.Select
                  value={language}
                  onChange={(e) => setLanguage(e.target.value)}
                >
                  {languageOptions.map((option) => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </Form.Select>
              </Form.Group>
            </Col>
            <Col md={6}>
              <Form.Group>
                <Form.Label>Number of Examples</Form.Label>
                <Form.Select
                  value={limit}
                  onChange={(e) => setLimit(parseInt(e.target.value))}
                >
                  <option value="3">3</option>
                  <option value="5">5</option>
                  <option value="10">10</option>
                </Form.Select>
              </Form.Group>
            </Col>
          </Row>

          <Button
            type="submit"
            variant="primary"
            size="lg"
            className="w-100"
            disabled={loading || !library.trim()}
          >
            {loading ? "Searching..." : "Explore Library Usage"}
          </Button>
        </Form>

        <div className="mt-4">
          <h5>Popular Libraries:</h5>
          <div className="d-flex flex-wrap gap-2">
            {getCurrentSuggestions().map((suggestion, index) => (
              <Button
                key={index}
                variant="outline-secondary"
                size="sm"
                onClick={() => handleSuggestionClick(suggestion)}
              >
                {suggestion}
              </Button>
            ))}
          </div>
        </div>
      </div>

      {loading && (
        <Spinner message="Searching GitHub for library usage examples..." />
      )}

      {error && (
        <Alert variant="danger" className="mb-4">
          {error}
        </Alert>
      )}

      {summary && (
        <Card className="shadow-sm mb-4">
          <Card.Header>
            <h4 className="mb-0">Library Usage Insights</h4>
          </Card.Header>
          <Card.Body>
            <div className="border rounded p-3 bg-light">
              <ReactMarkdown
                components={{
                  code({ node, inline, className, children, ...props }) {
                    const match = /language-(\w+)/.exec(className || "");
                    return !inline && match ? (
                      <SyntaxHighlighter
                        style={vscDarkPlus}
                        language={match[1]}
                        PreTag="div"
                        {...props}
                      >
                        {String(children).replace(/\n$/, "")}
                      </SyntaxHighlighter>
                    ) : (
                      <code className={className} {...props}>
                        {children}
                      </code>
                    );
                  },
                }}
              >
                {summary}
              </ReactMarkdown>
            </div>
          </Card.Body>
        </Card>
      )}

      {examples.length > 0 && (
        <div className="mb-4">
          <h2 className="mb-3">Usage Examples</h2>
          <div className="d-flex flex-column gap-4">
            {examples.map((example, index) => (
              <LibraryUsageExample
                key={index}
                example={example}
                index={index}
                language={language}
              />
            ))}
          </div>
        </div>
      )}
    </Container>
  );
}

export default LibraryExplorer;
