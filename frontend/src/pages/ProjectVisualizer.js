import React, { useState, useEffect, useRef } from "react";
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
import mermaid from "mermaid";

// Components
import Spinner from "../components/Spinner";

function ProjectVisualizer() {
  const [repoUrl, setRepoUrl] = useState("");
  const [folderPath, setFolderPath] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [diagramCode, setDiagramCode] = useState("");
  const [diagramType, setDiagramType] = useState("flowchart");
  const [activeTab, setActiveTab] = useState("github");

  const mermaidRef = useRef(null);

  const API_BASE_URL =
    process.env.REACT_APP_API_BASE_URL || "http://localhost:5000/api";

  useEffect(() => {
    // Initialize mermaid
    mermaid.initialize({
      startOnLoad: true,
      theme: "default",
    });
  }, []);

  useEffect(() => {
    // Render mermaid diagram when code changes
    if (diagramCode && mermaidRef.current) {
      try {
        mermaid.contentLoaded();
      } catch (e) {
        console.error("Mermaid rendering error:", e);
      }
    }
  }, [diagramCode]);

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
    setDiagramCode("");

    try {
      // Prepare request data
      const requestData = {};

      if (activeTab === "github") {
        requestData.repoUrl = repoUrl;
      } else {
        requestData.folderPath = folderPath;
      }

      // Make API call
      const response = await axios.post(
        `${API_BASE_URL}/repo/visualize-project`,
        requestData
      );

      // When receiving diagramCode from the API
      if (response.data && response.data.success) {
        // Clean the diagram code by removing triple backticks and "mermaid" declaration
        let cleanedCode = response.data.data.diagramCode
          .replace(/```mermaid\s*/g, "") // Remove opening ```mermaid
          .replace(/```\s*$/g, "") // Remove closing ```
          .trim();

        setDiagramCode(cleanedCode);
        setDiagramType(response.data.data.type || "flowchart");

        // Clear and re-render the diagram
        if (mermaidRef.current) {
          mermaidRef.current.innerHTML = "";
          mermaidRef.current.innerHTML = `<div class="mermaid">${cleanedCode}</div>`;
          mermaid.contentLoaded();
        }
      } else {
        throw new Error(
          response.data.error || "Failed to generate project visualization"
        );
      }
    } catch (err) {
      setError(
        err.message || "An error occurred while generating visualization"
      );
      console.error("Visualization error:", err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="project-visualizer">
      <Container className="py-4">
        <div className="text-center mb-4">
          <h1>Project Structure Visualizer</h1>
          <p className="lead">
            Generate a visual representation of your project's architecture and
            structure
          </p>
        </div>

        <Row>
          <Col lg={5} className="mb-4">
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

                  <Button
                    type="submit"
                    variant="primary"
                    className="w-100"
                    disabled={loading}
                  >
                    {loading ? "Generating..." : "Generate Visualization"}
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

          <Col lg={7}>
            <Card className="h-100">
              <Card.Header className="d-flex justify-content-between align-items-center">
                <h4 className="mb-0">Project Visualization</h4>
                {diagramCode && (
                  <Button
                    variant="outline-secondary"
                    size="sm"
                    onClick={() => {
                      const blob = new Blob([diagramCode], {
                        type: "text/plain",
                      });
                      const url = URL.createObjectURL(blob);
                      const a = document.createElement("a");
                      a.href = url;
                      a.download = "project-diagram.mmd";
                      document.body.appendChild(a);
                      a.click();
                      document.body.removeChild(a);
                      URL.revokeObjectURL(url);
                    }}
                  >
                    Download Diagram Code
                  </Button>
                )}
              </Card.Header>
              <Card.Body>
                {loading ? (
                  <Spinner message="Analyzing project structure and generating visualization..." />
                ) : diagramCode ? (
                  <Tab.Container defaultActiveKey="diagram">
                    <Nav variant="tabs" className="mb-3">
                      <Nav.Item>
                        <Nav.Link eventKey="diagram">Diagram</Nav.Link>
                      </Nav.Item>
                      <Nav.Item>
                        <Nav.Link eventKey="code">Mermaid Code</Nav.Link>
                      </Nav.Item>
                    </Nav>

                    <Tab.Content>
                      <Tab.Pane eventKey="diagram">
                        <div
                          className="diagram-container border rounded p-3 bg-light overflow-auto"
                          style={{ maxHeight: "600px" }}
                        >
                          <div ref={mermaidRef} className="text-center">
                            <div className="mermaid">{diagramCode}</div>
                          </div>
                        </div>
                      </Tab.Pane>
                      <Tab.Pane eventKey="code">
                        <div
                          className="code-container border rounded overflow-auto"
                          style={{ maxHeight: "600px" }}
                        >
                          <SyntaxHighlighter
                            language="markdown"
                            style={vscDarkPlus}
                            showLineNumbers
                          >
                            {diagramCode}
                          </SyntaxHighlighter>
                        </div>
                      </Tab.Pane>
                    </Tab.Content>
                  </Tab.Container>
                ) : (
                  <div className="text-center p-5 text-muted">
                    <p className="mb-0">
                      Your project visualization will appear here
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

export default ProjectVisualizer;
