import React, { useState } from "react";
import {
  Container,
  Row,
  Col,
  Form,
  Button,
  Card,
  Alert,
  Tab,
  Nav,
} from "react-bootstrap";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { vscDarkPlus } from "react-syntax-highlighter/dist/esm/styles/prism";
import axios from "axios";
import ReactMarkdown from "react-markdown";

// Components
import Spinner from "../components/Spinner";

function ReadmeGenerator() {
  const [repoUrl, setRepoUrl] = useState("");
  const [folderPath, setFolderPath] = useState("");
  const [projectName, setProjectName] = useState("");
  const [description, setDescription] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [readmeContent, setReadmeContent] = useState("");
  const [activeTab, setActiveTab] = useState("github");

  const API_BASE_URL =
    process.env.REACT_APP_API_BASE_URL || "http://localhost:5000/api";

  const handleSubmit = async (e) => {
    e.preventDefault();

    // Validate inputs
    if (
      (activeTab === "github" && !repoUrl.trim()) ||
      (activeTab === "local" && !folderPath.trim())
    ) {
      setError(
        activeTab === "github"
          ? "Please enter a GitHub repository URL"
          : "Please enter a local folder path"
      );
      return;
    }

    setLoading(true);
    setError(null);
    setReadmeContent("");

    try {
      // Prepare request data
      const requestData = {
        projectName,
        description,
      };

      if (activeTab === "github") {
        requestData.repoUrl = repoUrl;
      } else {
        requestData.folderPath = folderPath;
      }

      // Make API call
      const response = await axios.post(
        `${API_BASE_URL}/repo/generate-readme`,
        requestData
      );

      if (response.data && response.data.success) {
        setReadmeContent(response.data.data.content);
      } else {
        throw new Error(
          response.data.error || "Failed to generate README content"
        );
      }
    } catch (err) {
      setError(err.message || "An error occurred while generating README");
      console.error("README generation error:", err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="readme-generator">
      <Container className="py-4">
        <div className="text-center mb-4">
          <h1>README Generator</h1>
          <p className="lead">
            Generate a comprehensive README.md file for your project from its
            repository structure
          </p>
        </div>

        <Row>
          <Col lg={6} className="mb-4">
            <Card className="h-100">
              <Card.Header>
                <h4 className="mb-0">Project Information</h4>
              </Card.Header>
              <Card.Body>
                <Form onSubmit={handleSubmit}>
                  <Tab.Container
                    activeKey={activeTab}
                    onSelect={(k) => setActiveTab(k)}
                  >
                    <Nav variant="tabs" className="mb-3">
                      <Nav.Item>
                        <Nav.Link eventKey="github">GitHub Repository</Nav.Link>
                      </Nav.Item>
                      <Nav.Item>
                        <Nav.Link eventKey="local">Local Folder</Nav.Link>
                      </Nav.Item>
                    </Nav>

                    <Tab.Content>
                      <Tab.Pane eventKey="github">
                        <Form.Group className="mb-3">
                          <Form.Label>GitHub Repository URL</Form.Label>
                          <Form.Control
                            type="text"
                            placeholder="https://github.com/username/repo"
                            value={repoUrl}
                            onChange={(e) => setRepoUrl(e.target.value)}
                          />
                          <Form.Text className="text-muted">
                            Enter the full URL to a GitHub repository
                          </Form.Text>
                        </Form.Group>
                      </Tab.Pane>

                      <Tab.Pane eventKey="local">
                        <Form.Group className="mb-3">
                          <Form.Label>Local Folder Path</Form.Label>
                          <Form.Control
                            type="text"
                            placeholder="/path/to/your/project"
                            value={folderPath}
                            onChange={(e) => setFolderPath(e.target.value)}
                          />
                          <Form.Text className="text-muted">
                            Enter the absolute path to your project folder
                          </Form.Text>
                        </Form.Group>
                      </Tab.Pane>
                    </Tab.Content>
                  </Tab.Container>

                  <Form.Group className="mb-3">
                    <Form.Label>Project Name (Optional)</Form.Label>
                    <Form.Control
                      type="text"
                      placeholder="Project Name"
                      value={projectName}
                      onChange={(e) => setProjectName(e.target.value)}
                    />
                  </Form.Group>

                  <Form.Group className="mb-3">
                    <Form.Label>Project Description (Optional)</Form.Label>
                    <Form.Control
                      as="textarea"
                      rows={3}
                      placeholder="A brief description of your project"
                      value={description}
                      onChange={(e) => setDescription(e.target.value)}
                    />
                  </Form.Group>

                  <Button
                    type="submit"
                    variant="primary"
                    className="w-100"
                    disabled={loading}
                  >
                    {loading ? "Generating..." : "Generate README"}
                  </Button>
                </Form>

                {error && (
                  <Alert variant="danger" className="mt-3">
                    {error}
                  </Alert>
                )}
              </Card.Body>
            </Card>
          </Col>

          <Col lg={6}>
            <Card className="h-100">
              <Card.Header className="d-flex justify-content-between align-items-center">
                <h4 className="mb-0">Generated README</h4>
                {readmeContent && (
                  <Button
                    variant="outline-secondary"
                    size="sm"
                    onClick={() => {
                      const blob = new Blob([readmeContent], {
                        type: "text/markdown",
                      });
                      const url = URL.createObjectURL(blob);
                      const a = document.createElement("a");
                      a.href = url;
                      a.download = "README.md";
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
                  <Spinner message="Analyzing repository and generating README..." />
                ) : readmeContent ? (
                  <Tab.Container defaultActiveKey="preview">
                    <Nav variant="tabs" className="mb-3">
                      <Nav.Item>
                        <Nav.Link eventKey="preview">Preview</Nav.Link>
                      </Nav.Item>
                      <Nav.Item>
                        <Nav.Link eventKey="markdown">Markdown</Nav.Link>
                      </Nav.Item>
                    </Nav>

                    <Tab.Content>
                      <Tab.Pane eventKey="preview">
                        <div
                          className="readme-preview border rounded p-3 bg-light overflow-auto"
                          style={{ maxHeight: "600px" }}
                        >
                          <ReactMarkdown>{readmeContent}</ReactMarkdown>
                        </div>
                      </Tab.Pane>
                      <Tab.Pane eventKey="markdown">
                        <div
                          className="code-container border rounded overflow-auto"
                          style={{ maxHeight: "600px" }}
                        >
                          <SyntaxHighlighter
                            language="markdown"
                            style={vscDarkPlus}
                            showLineNumbers
                          >
                            {readmeContent}
                          </SyntaxHighlighter>
                        </div>
                      </Tab.Pane>
                    </Tab.Content>
                  </Tab.Container>
                ) : (
                  <div className="text-center p-5 text-muted">
                    <p className="mb-0">
                      Your generated README will appear here
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

export default ReadmeGenerator;
