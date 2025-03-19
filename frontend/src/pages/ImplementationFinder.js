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
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { vscDarkPlus } from "react-syntax-highlighter/dist/esm/styles/prism";
import ReactMarkdown from "react-markdown";

// Components
import Spinner from "../components/Spinner";

function ImplementationFinder() {
  const [functionName, setFunctionName] = useState("");
  const [language, setLanguage] = useState("javascript");
  const [limit, setLimit] = useState(3);
  const [implementations, setImplementations] = useState([]);
  const [comparison, setComparison] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [loadingComparison, setLoadingComparison] = useState(false);

  const API_BASE_URL =
    process.env.REACT_APP_API_BASE_URL || "http://localhost:5000/api";

  const languageOptions = [
    { value: "javascript", label: "JavaScript" },
    { value: "python", label: "Python" },
    { value: "java", label: "Java" },
    { value: "csharp", label: "C#" },
    { value: "typescript", label: "TypeScript" },
    { value: "go", label: "Go" },
    { value: "ruby", label: "Ruby" },
  ];

  const popularFunctions = {
    javascript: ["quicksort", "fibonacci", "debounce", "throttle", "mergesort"],
    python: [
      "quicksort",
      "fibonacci",
      "binary_search",
      "merge_sort",
      "breadth_first_search",
    ],
    java: [
      "quicksort",
      "fibonacci",
      "binarySearch",
      "mergeSort",
      "depthFirstSearch",
    ],
  };

  const getCurrentSuggestions = () => {
    return popularFunctions[language] || popularFunctions.javascript;
  };

  const handleSuggestionClick = (suggestion) => {
    setFunctionName(suggestion);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!functionName.trim()) {
      setError("Please enter a function or algorithm name");
      return;
    }

    try {
      setLoading(true);
      setError(null);
      setImplementations([]);
      setComparison(null);

      const response = await axios.get(
        `${API_BASE_URL}/github/implementations`,
        {
          params: {
            functionName,
            language,
            limit,
          },
        }
      );

      if (response.data && response.data.success) {
        setImplementations(response.data.data);

        if (response.data.data.length >= 2) {
          await generateComparison(response.data.data);
        }
      } else {
        throw new Error("Failed to fetch implementations");
      }
    } catch (err) {
      setError(err.message || "An error occurred");
    } finally {
      setLoading(false);
    }
  };

  const generateComparison = async (implementations) => {
    try {
      setLoadingComparison(true);

      const response = await axios.post(
        `${API_BASE_URL}/implementations/compare`,
        {
          implementations,
          language,
        }
      );

      if (response.data && response.data.success) {
        setComparison(response.data.data.comparison);
      }
    } catch (err) {
      console.error("Failed to generate comparison:", err);
    } finally {
      setLoadingComparison(false);
    }
  };

  return (
    <div className="implementation-finder">
      <Container className="py-4">
        <div className="text-center mb-4">
          <h1>Implementation Finder</h1>
          <p className="lead">
            Search for different implementations of algorithms or functions
            across GitHub repositories to compare approaches.
          </p>
        </div>

        <Card className="mb-4">
          <Card.Body>
            <Form onSubmit={handleSubmit}>
              <Form.Group className="mb-3">
                <Form.Label>Function or Algorithm Name</Form.Label>
                <Form.Control
                  type="text"
                  placeholder="Enter a function or algorithm name (e.g., quicksort)"
                  value={functionName}
                  onChange={(e) => setFunctionName(e.target.value)}
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
                    <Form.Label>Number of Implementations</Form.Label>
                    <Form.Select
                      value={limit}
                      onChange={(e) => setLimit(parseInt(e.target.value))}
                    >
                      <option value="2">2</option>
                      <option value="3">3</option>
                      <option value="5">5</option>
                    </Form.Select>
                  </Form.Group>
                </Col>
              </Row>

              <Button
                type="submit"
                variant="primary"
                className="w-100"
                disabled={loading || !functionName.trim()}
              >
                {loading ? "Searching..." : "Find Implementations"}
              </Button>
            </Form>

            <div className="mt-4">
              <h5>Popular Functions:</h5>
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
          </Card.Body>
        </Card>

        {loading && (
          <Spinner message="Searching GitHub for implementations..." />
        )}

        {error && (
          <Alert variant="danger" className="mb-4">
            {error}
          </Alert>
        )}

        {loadingComparison && (
          <Spinner message="Generating implementation comparison..." />
        )}

        {comparison && (
          <Card className="mb-4">
            <Card.Header>
              <h4 className="mb-0">Implementation Comparison</h4>
            </Card.Header>
            <Card.Body>
              <div className="border rounded p-3 bg-light">
                <ReactMarkdown>{comparison}</ReactMarkdown>
              </div>
            </Card.Body>
          </Card>
        )}

        {implementations.length > 0 && (
          <div>
            <h2 className="mb-3">Found Implementations</h2>
            {implementations.map((implementation, index) => (
              <Card key={index} className="mb-4">
                <Card.Header>
                  <h5 className="mb-0">Implementation {index + 1}</h5>
                </Card.Header>
                <Card.Body>
                  <div>
                    <p className="mb-3">
                      <strong>Repository:</strong>{" "}
                      <a
                        href={implementation.repository.url}
                        target="_blank"
                        rel="noopener noreferrer"
                      >
                        {implementation.repository.full_name}
                      </a>
                    </p>
                    <p className="mb-3">
                      <strong>File:</strong>{" "}
                      <a
                        href={implementation.file.url}
                        target="_blank"
                        rel="noopener noreferrer"
                      >
                        {implementation.file.path}
                      </a>
                    </p>
                    <div className="border rounded">
                      <SyntaxHighlighter
                        language={language}
                        style={vscDarkPlus}
                      >
                        {implementation.file.content}
                      </SyntaxHighlighter>
                    </div>
                  </div>
                </Card.Body>
              </Card>
            ))}
          </div>
        )}
      </Container>
    </div>
  );
}

export default ImplementationFinder;
