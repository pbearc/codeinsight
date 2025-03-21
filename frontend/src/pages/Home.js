import React from "react";
import { Container, Row, Col, Button, Card } from "react-bootstrap";
import { Link } from "react-router-dom";

function Home() {
  return (
    <Container className="py-5">
      {/* Hero Section */}
      <div className="text-center mb-5">
        <h1 className="display-4 mb-3">Welcome to CodeInsight</h1>
        <p className="lead mb-4">
          Your intelligent assistant for understanding code, exploring
          libraries, generating documentation, and analyzing developers.
        </p>
        <div className="d-flex justify-content-center gap-3">
          <Button as={Link} to="/code-analyzer" variant="primary" size="lg">
            Analyze Code
          </Button>
          <Button
            as={Link}
            to="/developer-analysis"
            variant="outline-secondary"
            size="lg"
          >
            Analyze Developers
          </Button>
        </div>
      </div>

      {/* Features Section */}
      <div className="mb-5">
        <h2 className="text-center mb-4">Features</h2>
        <Row className="g-4">
          <Col md={6} lg={3}>
            <Card className="h-100 shadow-sm">
              <Card.Body>
                <h3 className="h5 mb-3">Code Analyzer</h3>
                <p className="card-text mb-3">
                  Upload code snippets to identify patterns, potential bugs, and
                  get improvement suggestions.
                </p>
                <Button
                  as={Link}
                  to="/code-analyzer"
                  variant="outline-primary"
                  size="sm"
                >
                  Try it
                </Button>
              </Card.Body>
            </Card>
          </Col>

          <Col md={6} lg={3}>
            <Card className="h-100 shadow-sm">
              <Card.Body>
                <h3 className="h5 mb-3">Library Explorer</h3>
                <p className="card-text mb-3">
                  Find common usage patterns and examples for any library or
                  framework.
                </p>
                <Button
                  as={Link}
                  to="/library-explorer"
                  variant="outline-primary"
                  size="sm"
                >
                  Try it
                </Button>
              </Card.Body>
            </Card>
          </Col>

          <Col md={6} lg={3}>
            <Card className="h-100 shadow-sm">
              <Card.Body>
                <h3 className="h5 mb-3">Documentation Generator</h3>
                <p className="card-text mb-3">
                  Transform complex code into clear, concise documentation with
                  inline comments or documentation files.
                </p>
                <Button
                  as={Link}
                  to="/documentation-generator"
                  variant="outline-primary"
                  size="sm"
                >
                  Try it
                </Button>
              </Card.Body>
            </Card>
          </Col>

          <Col md={6} lg={3}>
            <Card className="h-100 shadow-sm">
              <Card.Body>
                <h3 className="h5 mb-3">Implementation Finder</h3>
                <p className="card-text mb-3">
                  Compare different implementations of algorithms or functions
                  across GitHub repositories.
                </p>
                <Button
                  as={Link}
                  to="/implementation-finder"
                  variant="outline-primary"
                  size="sm"
                >
                  Try it
                </Button>
              </Card.Body>
            </Card>
          </Col>
        </Row>

        {/* New Row for Additional Features */}
        <Row className="g-4 mt-2">
          <Col md={6} lg={4}>
            <Card className="h-100 shadow-sm">
              <Card.Body>
                <h3 className="h5 mb-3">README Generator</h3>
                <p className="card-text mb-3">
                  Automatically generate professional README.md files for your
                  projects by analyzing your repository structure and code.
                </p>
                <Button
                  as={Link}
                  to="/readme-generator"
                  variant="outline-primary"
                  size="sm"
                >
                  Try it
                </Button>
              </Card.Body>
            </Card>
          </Col>

          <Col md={6} lg={4}>
            <Card className="h-100 shadow-sm">
              <Card.Body>
                <h3 className="h5 mb-3">Project Visualizer</h3>
                <p className="card-text mb-3">
                  Create visual diagrams of your project's architecture and
                  structure to better understand and document your codebase.
                </p>
                <Button
                  as={Link}
                  to="/project-visualizer"
                  variant="outline-primary"
                  size="sm"
                >
                  Try it
                </Button>
              </Card.Body>
            </Card>
          </Col>

          {/* New Developer Analysis Feature Card */}
          <Col md={6} lg={4}>
            <Card className="h-100 shadow-sm border-primary">
              <Card.Body>
                <span className="badge bg-primary mb-2">New</span>
                <h3 className="h5 mb-3">Developer Analysis</h3>
                <p className="card-text mb-3">
                  Analyze GitHub developer profiles to understand their
                  technical skills, project complexity, and evaluate potential
                  team members.
                </p>
                <div className="d-flex gap-2">
                  <Button
                    as={Link}
                    to="/developer-analysis"
                    variant="outline-primary"
                    size="sm"
                  >
                    Analyze Dev
                  </Button>
                  <Button
                    as={Link}
                    to="/developer-comparison"
                    variant="outline-primary"
                    size="sm"
                  >
                    Compare Devs
                  </Button>
                </div>
              </Card.Body>
            </Card>
          </Col>
        </Row>
      </div>

      {/* How It Works Section */}
      <div>
        <h2 className="text-center mb-4">How It Works</h2>
        <Row className="g-4 text-center">
          <Col md={3}>
            <div className="d-flex flex-column align-items-center">
              <div
                className="rounded-circle bg-primary text-white d-flex align-items-center justify-content-center mb-3"
                style={{ width: "60px", height: "60px" }}
              >
                <h3 className="mb-0">1</h3>
              </div>
              <h4>Input Your Code or Query</h4>
              <p>
                Enter a code snippet, GitHub URL, or developer username that you
                want to analyze.
              </p>
            </div>
          </Col>

          <Col md={3}>
            <div className="d-flex flex-column align-items-center">
              <div
                className="rounded-circle bg-primary text-white d-flex align-items-center justify-content-center mb-3"
                style={{ width: "60px", height: "60px" }}
              >
                <h3 className="mb-0">2</h3>
              </div>
              <h4>Search GitHub</h4>
              <p>
                Our system searches through GitHub repositories to find relevant
                code, examples, and developer activity.
              </p>
            </div>
          </Col>

          <Col md={3}>
            <div className="d-flex flex-column align-items-center">
              <div
                className="rounded-circle bg-primary text-white d-flex align-items-center justify-content-center mb-3"
                style={{ width: "60px", height: "60px" }}
              >
                <h3 className="mb-0">3</h3>
              </div>
              <h4>AI Analysis</h4>
              <p>
                Our LLM analyzes the data and generates insights, explanations,
                documentation, and technical profiles.
              </p>
            </div>
          </Col>

          <Col md={3}>
            <div className="d-flex flex-column align-items-center">
              <div
                className="rounded-circle bg-primary text-white d-flex align-items-center justify-content-center mb-3"
                style={{ width: "60px", height: "60px" }}
              >
                <h3 className="mb-0">4</h3>
              </div>
              <h4>View Results</h4>
              <p>
                Review the generated insights, examples, documentation, or
                developer profiles to make better decisions.
              </p>
            </div>
          </Col>
        </Row>
      </div>
    </Container>
  );
}

export default Home;
