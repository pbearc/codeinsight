// frontend/src/pages/CodeAnalyzer.js
import React, { useState } from "react";
import axios from "axios";
import {
  Container,
  Row,
  Col,
  Form,
  Button,
  Alert,
  Nav,
  Tab,
} from "react-bootstrap";

// Components
import Spinner from "../components/Spinner";
import CodeEditor from "../components/CodeEditor";
import AnalysisResult from "../components/AnalysisResult";

function CodeAnalyzer() {
  const [code, setCode] = useState("");
  const [githubUrl, setGithubUrl] = useState("");
  const [language, setLanguage] = useState("javascript");
  const [analysisType, setAnalysisType] = useState("analyze");
  const [result, setResult] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [activeTab, setActiveTab] = useState("code"); // 'code' or 'github'

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
    { value: "swift", label: "Swift" },
    { value: "kotlin", label: "Kotlin" },
    { value: "rust", label: "Rust" },
  ];

  const analysisTypeOptions = [
    { value: "analyze", label: "Comprehensive Analysis" },
    { value: "explain", label: "Simple Explanation" },
    { value: "patternIdentification", label: "Pattern Identification" },
  ];

  // Function to extract owner, repo, and path from GitHub URL
  const parseGitHubUrl = (url) => {
    try {
      // Handle raw.githubusercontent.com URLs
      if (url.includes("raw.githubusercontent.com")) {
        const parts = url
          .replace("https://raw.githubusercontent.com/", "")
          .split("/");
        const owner = parts[0];
        const repo = parts[1];
        const path = parts.slice(3).join("/");
        return { owner, repo, path };
      }

      // Handle github.com URLs
      if (url.includes("github.com")) {
        const parts = url.replace("https://github.com/", "").split("/");
        const owner = parts[0];
        const repo = parts[1];

        // Handle blob URLs
        if (parts[2] === "blob") {
          const path = parts.slice(4).join("/");
          return { owner, repo, path };
        }
      }

      throw new Error("Invalid GitHub URL format");
    } catch (err) {
      throw new Error("Failed to parse GitHub URL: " + err.message);
    }
  };

  // Function to fetch code from GitHub
  const fetchCodeFromGitHub = async (url) => {
    try {
      const { owner, repo, path } = parseGitHubUrl(url);

      const response = await axios.get(`${API_BASE_URL}/github/content`, {
        params: { owner, repo, path },
      });

      if (response.data && response.data.success) {
        return response.data.data;
      } else {
        throw new Error("Failed to fetch code from GitHub");
      }
    } catch (err) {
      throw new Error("Failed to fetch code: " + err.message);
    }
  };

  const handleAnalyze = async () => {
    try {
      setLoading(true);
      setError(null);
      setResult(null);

      let codeToAnalyze = code;

      // If using GitHub URL input, fetch the code first
      if (activeTab === "github" && githubUrl) {
        codeToAnalyze = await fetchCodeFromGitHub(githubUrl);
        setCode(codeToAnalyze); // Update the code editor with fetched code
      }

      if (!codeToAnalyze.trim()) {
        throw new Error("No code to analyze");
      }

      // Analyze the code
      const response = await axios.post(`${API_BASE_URL}/analyze`, {
        code: codeToAnalyze,
        language,
        analysisType,
      });

      if (response.data && response.data.success) {
        setResult(response.data.data);
      } else {
        throw new Error("Analysis failed");
      }
    } catch (err) {
      setError(err.message || "An error occurred");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container className="py-4">
      <div className="text-center mb-4">
        <h1>Code Analyzer</h1>
        <p className="lead">
          Upload code snippets or provide GitHub URLs to analyze code patterns,
          identify issues, and get improvement suggestions.
        </p>
      </div>

      <div className="bg-white p-4 border rounded shadow-sm mb-4">
        <Tab.Container activeKey={activeTab} onSelect={(k) => setActiveTab(k)}>
          <Nav variant="tabs" className="mb-3">
            <Nav.Item>
              <Nav.Link eventKey="code">Enter Code</Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey="github">GitHub URL</Nav.Link>
            </Nav.Item>
          </Nav>
          <Tab.Content>
            <Tab.Pane eventKey="code">
              <Form.Group className="mb-3">
                <Form.Label>Enter Code</Form.Label>
                <CodeEditor
                  value={code}
                  onChange={setCode}
                  language={language}
                />
              </Form.Group>
            </Tab.Pane>
            <Tab.Pane eventKey="github">
              <Form.Group className="mb-3">
                <Form.Label>
                  GitHub URL (to raw file or repository file)
                </Form.Label>
                <Form.Control
                  type="text"
                  placeholder="https://github.com/username/repo/blob/main/path/to/file.js"
                  value={githubUrl}
                  onChange={(e) => setGithubUrl(e.target.value)}
                />
              </Form.Group>
            </Tab.Pane>
          </Tab.Content>
        </Tab.Container>

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
              <Form.Label>Analysis Type</Form.Label>
              <Form.Select
                value={analysisType}
                onChange={(e) => setAnalysisType(e.target.value)}
              >
                {analysisTypeOptions.map((option) => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </Form.Select>
            </Form.Group>
          </Col>
        </Row>

        <Button
          variant="primary"
          size="lg"
          className="w-100"
          onClick={handleAnalyze}
          disabled={
            loading ||
            (activeTab === "code" && !code.trim()) ||
            (activeTab === "github" && !githubUrl.trim())
          }
        >
          {loading ? "Analyzing..." : "Analyze Code"}
        </Button>
      </div>

      {loading && <Spinner message="Analyzing your code..." />}

      {error && (
        <Alert variant="danger" className="mb-4">
          {error}
        </Alert>
      )}

      {result && <AnalysisResult result={result} />}
    </Container>
  );
}

export default CodeAnalyzer;
