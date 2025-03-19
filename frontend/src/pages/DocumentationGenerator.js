// frontend/src/pages/DocumentationGenerator.js
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
  Card,
} from "react-bootstrap";
import ReactMarkdown from "react-markdown";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { vscDarkPlus } from "react-syntax-highlighter/dist/esm/styles/prism";

// Components
import Spinner from "../components/Spinner";
import CodeEditor from "../components/CodeEditor";

function DocumentationGenerator() {
  const [code, setCode] = useState("");
  const [githubUrl, setGithubUrl] = useState("");
  const [language, setLanguage] = useState("javascript");
  const [documentation, setDocumentation] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [activeTab, setActiveTab] = useState("code"); // 'code' or 'github'
  const [useExamples, setUseExamples] = useState(true);

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

  const handleGenerate = async () => {
    try {
      setLoading(true);
      setError(null);
      setDocumentation(null);

      let codeToDocument = code;

      // If using GitHub URL input, fetch the code first
      if (activeTab === "github" && githubUrl) {
        codeToDocument = await fetchCodeFromGitHub(githubUrl);
        setCode(codeToDocument); // Update the code editor with fetched code
      }

      if (!codeToDocument.trim()) {
        throw new Error("No code to document");
      }

      // Generate documentation
      const response = await axios.post(`${API_BASE_URL}/generate-docs`, {
        code: codeToDocument,
        language,
        useExamples,
      });

      if (response.data && response.data.success) {
        setDocumentation(response.data.data.documentation);
      } else {
        throw new Error("Documentation generation failed");
      }
    } catch (err) {
      setError(err.message || "An error occurred");
    } finally {
      setLoading(false);
    }
  };

  const handleCopyToClipboard = () => {
    if (documentation) {
      navigator.clipboard
        .writeText(documentation)
        .then(() => {
          alert("Documentation copied to clipboard!");
        })
        .catch((err) => {
          console.error("Failed to copy: ", err);
        });
    }
  };

  const handleDownload = () => {
    if (documentation) {
      const element = document.createElement("a");
      const file = new Blob([documentation], { type: "text/markdown" });
      element.href = URL.createObjectURL(file);
      element.download = `documentation_${new Date()
        .toISOString()
        .slice(0, 10)}.md`;
      document.body.appendChild(element);
      element.click();
      document.body.removeChild(element);
    }
  };

  return (
    <Container className="py-4">
      <div className="text-center mb-4">
        <h1>Documentation Generator</h1>
        <p className="lead">
          Transform complex code into clear, concise documentation with
          real-world examples.
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
          <Col md={6} className="d-flex align-items-end">
            <Form.Check
              type="checkbox"
              id="use-examples-checkbox"
              label="Include real-world examples"
              checked={useExamples}
              onChange={(e) => setUseExamples(e.target.checked)}
              className="mb-3"
            />
          </Col>
        </Row>

        <Button
          variant="primary"
          size="lg"
          className="w-100"
          onClick={handleGenerate}
          disabled={
            loading ||
            (activeTab === "code" && !code.trim()) ||
            (activeTab === "github" && !githubUrl.trim())
          }
        >
          {loading ? "Generating..." : "Generate Documentation"}
        </Button>
      </div>

      {loading && <Spinner message="Generating documentation..." />}

      {error && (
        <Alert variant="danger" className="mb-4">
          {error}
        </Alert>
      )}

      {documentation && (
        <Card className="shadow-sm mb-4">
          <Card.Header>
            <div className="d-flex justify-content-between align-items-center">
              <h4 className="mb-0">Generated Documentation</h4>
              <div>
                <Button
                  variant="outline-secondary"
                  className="me-2"
                  onClick={handleCopyToClipboard}
                >
                  Copy to Clipboard
                </Button>
                <Button variant="outline-secondary" onClick={handleDownload}>
                  Download as Markdown
                </Button>
              </div>
            </div>
          </Card.Header>
          <Card.Body>
            <div className="documentation-content border rounded p-3 bg-light">
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
                {documentation}
              </ReactMarkdown>
            </div>
          </Card.Body>
        </Card>
      )}
    </Container>
  );
}

export default DocumentationGenerator;
