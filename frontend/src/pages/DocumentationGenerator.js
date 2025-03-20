import React, { useState } from "react";
import {
  Container,
  Row,
  Col,
  Form,
  Button,
  Alert,
  Card,
  Tab,
  Nav,
} from "react-bootstrap";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { vscDarkPlus } from "react-syntax-highlighter/dist/esm/styles/prism";
import axios from "axios";
import ReactMarkdown from "react-markdown";

// Components
import Spinner from "../components/Spinner";

function DocumentationGenerator() {
  const [code, setCode] = useState("");
  const [githubUrl, setGithubUrl] = useState("");
  const [language, setLanguage] = useState("javascript");
  const [useExamples, setUseExamples] = useState(false);
  const [inlineDocumentation, setInlineDocumentation] = useState(false);
  const [activeTab, setActiveTab] = useState("code");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [documentation, setDocumentation] = useState("");

  const API_BASE_URL =
    process.env.REACT_APP_API_BASE_URL || "http://localhost:5000/api";

  const languageOptions = [
    { value: "javascript", label: "JavaScript" },
    { value: "python", label: "Python" },
    { value: "java", label: "Java" },
    { value: "go", label: "Go" },
    { value: "ruby", label: "Ruby" },
    { value: "php", label: "PHP" },
    { value: "typescript", label: "TypeScript" },
    { value: "csharp", label: "C#" },
  ];

  const fetchGithubContent = async () => {
    if (!githubUrl.trim()) {
      setError("Please enter a GitHub URL");
      return;
    }

    setLoading(true);
    setError(null);

    try {
      // Extract owner, repo, and path from GitHub URL
      const githubRegex = /github\.com\/([^\/]+)\/([^\/]+)\/blob\/[^\/]+\/(.+)/;
      const match = githubUrl.match(githubRegex);

      if (!match || match.length !== 4) {
        throw new Error(
          "Invalid GitHub URL. Please use a URL to a specific file."
        );
      }

      const [, owner, repo, path] = match;

      // Fetch file content from GitHub
      const response = await axios.get(
        `${API_BASE_URL}/github/content?owner=${owner}&repo=${repo}&path=${path}`
      );

      if (response.data && response.data.success) {
        setCode(response.data.data);

        // Try to determine language from file extension
        const fileExtension = path.split(".").pop().toLowerCase();
        const extensionMap = {
          js: "javascript",
          py: "python",
          java: "java",
          go: "go",
          rb: "ruby",
          php: "php",
          ts: "typescript",
          cs: "csharp",
        };

        if (extensionMap[fileExtension]) {
          setLanguage(extensionMap[fileExtension]);
        }

        setActiveTab("code");
      } else {
        throw new Error("Failed to fetch code from GitHub");
      }
    } catch (err) {
      setError(err.message || "An error occurred");
    } finally {
      setLoading(false);
    }
  };

  const handleGenerate = async () => {
    if (!code.trim()) {
      setError("Please enter or fetch some code first");
      return;
    }

    setLoading(true);
    setError(null);
    setDocumentation("");

    try {
      const response = await axios.post(`${API_BASE_URL}/generate-docs`, {
        code,
        language,
        useExamples,
        inline: inlineDocumentation,
      });

      if (response.data && response.data.success) {
        setDocumentation(response.data.data.documentation);
      } else {
        throw new Error("Failed to generate documentation");
      }
    } catch (err) {
      setError(err.message || "An error occurred");
      console.error("Documentation generation error:", err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="documentation-generator">
      <Container className="py-4">
        <div className="text-center mb-4">
          <h1>Documentation Generator</h1>
          <p className="lead">
            Transform complex code into clear, concise documentation
          </p>
        </div>

        <Row>
          <Col lg={6} className="mb-4">
            <Card>
              <Card.Header>
                <Nav
                  variant="tabs"
                  activeKey={activeTab}
                  onSelect={setActiveTab}
                >
                  <Nav.Item>
                    <Nav.Link eventKey="code">Enter Code</Nav.Link>
                  </Nav.Item>
                  <Nav.Item>
                    <Nav.Link eventKey="github">GitHub URL</Nav.Link>
                  </Nav.Item>
                </Nav>
              </Card.Header>
              <Card.Body>
                <Tab.Content>
                  <Tab.Pane eventKey="code">
                    <Form.Group className="mb-3">
                      <Form.Label>GitHub URL</Form.Label>
                      <Form.Control
                        type="text"
                        value={githubUrl}
                        onChange={(e) => setGithubUrl(e.target.value)}
                        placeholder="https://github.com/user/repo/blob/main/file.js"
                      />
                      <Form.Text className="text-muted">
                        Enter the URL to a file in a GitHub repository
                      </Form.Text>
                    </Form.Group>
                    <Button
                      variant="secondary"
                      onClick={fetchGithubContent}
                      disabled={loading}
                    >
                      Fetch Code
                    </Button>
                  </Tab.Pane>
                </Tab.Content>

                <div className="mt-3">
                  <Form.Group className="mb-3">
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

                  <Form.Group className="mb-3">
                    <Form.Check
                      type="checkbox"
                      label="Include real-world examples"
                      checked={useExamples}
                      onChange={(e) => setUseExamples(e.target.checked)}
                    />
                  </Form.Group>

                  <Form.Group className="mb-3">
                    <Form.Check
                      type="checkbox"
                      label="Generate inline documentation/comments"
                      checked={inlineDocumentation}
                      onChange={(e) => setInlineDocumentation(e.target.checked)}
                      id="inline-docs-checkbox"
                    />
                    <Form.Text className="text-muted">
                      When enabled, adds comments directly to the code. When
                      disabled, generates a separate documentation file.
                    </Form.Text>
                  </Form.Group>

                  <Button
                    variant="primary"
                    onClick={handleGenerate}
                    className="w-100"
                    disabled={loading || !code.trim()}
                  >
                    {loading ? "Generating..." : "Generate Documentation"}
                  </Button>
                </div>

                {error && (
                  <Alert variant="danger" className="mt-3">
                    {error}
                  </Alert>
                )}
              </Card.Body>
            </Card>
          </Col>

          <Col lg={6}>
            <Card>
              <Card.Header className="d-flex justify-content-between align-items-center">
                <h4 className="mb-0">Generated Documentation</h4>
                {documentation && (
                  <Button
                    variant="outline-secondary"
                    size="sm"
                    onClick={() => {
                      const blob = new Blob([documentation], {
                        type: inlineDocumentation
                          ? "text/plain"
                          : "text/markdown",
                      });
                      const url = URL.createObjectURL(blob);
                      const a = document.createElement("a");
                      a.href = url;
                      a.download = inlineDocumentation
                        ? `documented.${
                            language === "javascript" ? "js" : language
                          }`
                        : "documentation.md";
                      document.body.appendChild(a);
                      a.click();
                      document.body.removeChild(a);
                      URL.revokeObjectURL(url);
                    }}
                  >
                    Download
                  </Button>
                )}
              </Card.Header>
              <Card.Body>
                {loading ? (
                  <Spinner message="Generating documentation..." />
                ) : documentation ? (
                  inlineDocumentation ? (
                    <div
                      className="code-container border rounded overflow-auto"
                      style={{ maxHeight: "600px" }}
                    >
                      <SyntaxHighlighter
                        language={language}
                        style={vscDarkPlus}
                        showLineNumbers
                      >
                        {documentation}
                      </SyntaxHighlighter>
                    </div>
                  ) : (
                    <div
                      className="markdown-container border rounded p-3 bg-light overflow-auto"
                      style={{ maxHeight: "600px" }}
                    >
                      <ReactMarkdown>{documentation}</ReactMarkdown>
                    </div>
                  )
                ) : (
                  <div className="text-center p-5 text-muted">
                    <p className="mb-0">
                      Your generated documentation will appear here
                    </p>
                  </div>
                )}
              </Card.Body>
            </Card>
          </Col>
        </Row>
      </Container>
    </div>
  );
}

export default DocumentationGenerator;
