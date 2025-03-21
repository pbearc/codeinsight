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
  Badge,
  ListGroup,
  ProgressBar,
  Tabs,
  Tab,
} from "react-bootstrap";
import { XCircle } from "react-bootstrap-icons";

// Components
import Spinner from "../components/Spinner";

function DeveloperComparison() {
  const [usernames, setUsernames] = useState([""]);
  const [currentUsername, setCurrentUsername] = useState("");
  const [focus, setFocus] = useState("general");
  const [comparison, setComparison] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const API_BASE_URL =
    process.env.REACT_APP_API_BASE_URL || "http://localhost:5000/api";

  const focusOptions = [
    { value: "general", label: "General Comparison" },
    { value: "languages", label: "Programming Languages" },
    { value: "projects", label: "Project Complexity" },
    { value: "contributions", label: "Contribution Patterns" },
  ];

  const popularPairs = [
    ["torvalds", "gaearon"],
    ["yyx990803", "sindresorhus"],
    ["tj", "kentcdodds"],
  ];

  const handleAddUsername = (e) => {
    e.preventDefault();
    if (currentUsername.trim()) {
      setUsernames([...usernames, currentUsername.trim()]);
      setCurrentUsername("");
    }
  };

  const handleRemoveUsername = (index) => {
    const newUsernames = [...usernames];
    newUsernames.splice(index, 1);
    setUsernames(newUsernames);
  };

  const handleUsernameChange = (e, index) => {
    const newUsernames = [...usernames];
    newUsernames[index] = e.target.value;
    setUsernames(newUsernames);
  };

  const handleSuggestionClick = (pair) => {
    setUsernames(pair);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    // Filter out empty usernames
    const filteredUsernames = usernames.filter(
      (username) => username.trim() !== ""
    );

    if (filteredUsernames.length < 2) {
      setError("Please enter at least two GitHub usernames");
      return;
    }

    try {
      setLoading(true);
      setError(null);
      setComparison(null);

      const response = await axios.post(`${API_BASE_URL}/developers/compare`, {
        usernames: filteredUsernames,
        focus,
      });

      if (response.data && response.data.success) {
        setComparison(response.data.data.analysis);
      } else {
        throw new Error("Failed to compare developer profiles");
      }
    } catch (err) {
      setError(
        err.response?.data?.error ||
          err.message ||
          "An error occurred during comparison"
      );
    } finally {
      setLoading(false);
    }
  };

  const renderSkillComparison = () => {
    if (!comparison || !comparison.skillComparison) {
      return <p>No skill comparison data available</p>;
    }

    return (
      <div>
        {Object.entries(comparison.skillComparison).map(
          ([skill, scores], index) => (
            <div key={index} className="mb-4">
              <h5>{skill}</h5>
              {scores.map((score, developerIndex) => (
                <div key={developerIndex} className="mb-2">
                  <div className="d-flex justify-content-between mb-1">
                    <div>{comparison.developers[developerIndex]}</div>
                    <div>{score.toFixed(1)}/100</div>
                  </div>
                  <ProgressBar
                    now={score}
                    variant={getColorForIndex(developerIndex)}
                    style={{ height: "10px" }}
                  />
                </div>
              ))}
            </div>
          )
        )}
      </div>
    );
  };

  const renderProjectScaleComparison = () => {
    if (!comparison || !comparison.projectScale) {
      return <p>No project scale data available</p>;
    }

    return (
      <div>
        {Object.entries(comparison.projectScale).map(
          ([developer, scales], index) => {
            // Count projects by complexity
            const counts = {
              toy: 0,
              small: 0,
              medium: 0,
              large: 0,
              enterprise: 0,
            };

            scales.forEach((scale) => {
              const lowerScale = scale.toLowerCase();
              if (counts.hasOwnProperty(lowerScale)) {
                counts[lowerScale]++;
              }
            });

            const total = Object.values(counts).reduce(
              (sum, count) => sum + count,
              0
            );

            return (
              <div key={index} className="mb-4">
                <h5>{developer}</h5>
                <div className="d-flex mb-2">
                  {Object.entries(counts).map(([scale, count], i) => (
                    <div
                      key={i}
                      className="text-center"
                      style={{
                        flex: `${count / total}`,
                        backgroundColor: getColorForScale(scale),
                        color: "white",
                        padding: "5px 0",
                        fontSize: "0.8rem",
                      }}
                    >
                      {count > 0 ? `${scale} (${count})` : ""}
                    </div>
                  ))}
                </div>
              </div>
            );
          }
        )}
      </div>
    );
  };

  const renderComplementarityChart = () => {
    if (!comparison || !comparison.complementarity) {
      return <p>No complementarity data available</p>;
    }

    return (
      <div>
        {Object.entries(comparison.complementarity).map(
          ([skill, devs], index) => (
            <div key={index} className="mb-3">
              <h5>{skill}</h5>
              <div className="d-flex flex-wrap gap-2 mb-2">
                {devs.map((developer, devIndex) => (
                  <Badge
                    key={devIndex}
                    bg={getColorForIndex(
                      comparison.developers.indexOf(developer),
                      true
                    )}
                    className="p-2"
                  >
                    {developer}
                  </Badge>
                ))}
              </div>
            </div>
          )
        )}
      </div>
    );
  };

  const getColorForScale = (scale) => {
    const colors = {
      toy: "#17a2b8", // info
      small: "#28a745", // success
      medium: "#007bff", // primary
      large: "#ffc107", // warning
      enterprise: "#dc3545", // danger
    };
    return colors[scale.toLowerCase()] || "#6c757d"; // default to secondary
  };

  const getColorForIndex = (index, asBadge = false) => {
    const colors = [
      "primary",
      "success",
      "info",
      "warning",
      "danger",
      "secondary",
    ];
    return colors[index % colors.length];
  };

  return (
    <div className="developer-comparison">
      <Container className="py-4">
        <div className="text-center mb-4">
          <h1>Developer Comparison</h1>
          <p className="lead">
            Compare multiple GitHub developers to understand their relative
            strengths, specializations, and how they might work together.
          </p>
        </div>

        <Card className="mb-4">
          <Card.Body>
            <Form onSubmit={handleSubmit}>
              <Form.Group className="mb-3">
                <Form.Label>GitHub Usernames</Form.Label>
                {usernames.map((username, index) => (
                  <div key={index} className="d-flex mb-2">
                    <Form.Control
                      type="text"
                      placeholder={`GitHub username ${index + 1}`}
                      value={username}
                      onChange={(e) => handleUsernameChange(e, index)}
                      required={index < 2} // At least 2 usernames required
                    />
                    {index > 0 && (
                      <Button
                        variant="outline-danger"
                        className="ms-2"
                        onClick={() => handleRemoveUsername(index)}
                      >
                        <XCircle />
                      </Button>
                    )}
                  </div>
                ))}

                <div className="d-flex mt-2">
                  <Form.Control
                    type="text"
                    placeholder="Add another username"
                    value={currentUsername}
                    onChange={(e) => setCurrentUsername(e.target.value)}
                  />
                  <Button
                    variant="outline-secondary"
                    className="ms-2"
                    onClick={handleAddUsername}
                    disabled={!currentUsername.trim()}
                  >
                    Add
                  </Button>
                </div>
              </Form.Group>

              <Form.Group className="mb-3">
                <Form.Label>Comparison Focus</Form.Label>
                <Form.Select
                  value={focus}
                  onChange={(e) => setFocus(e.target.value)}
                >
                  {focusOptions.map((option) => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </Form.Select>
                <Form.Text className="text-muted">
                  Select what aspect of the developers you want to focus on
                </Form.Text>
              </Form.Group>

              <Button
                type="submit"
                variant="primary"
                className="w-100"
                disabled={
                  loading || usernames.filter((u) => u.trim()).length < 2
                }
              >
                {loading ? "Comparing..." : "Compare Developers"}
              </Button>
            </Form>

            <div className="mt-4">
              <h5>Popular Comparisons:</h5>
              <div className="d-flex flex-wrap gap-2">
                {popularPairs.map((pair, index) => (
                  <Button
                    key={index}
                    variant="outline-secondary"
                    size="sm"
                    onClick={() => handleSuggestionClick(pair)}
                  >
                    {pair.join(" vs ")}
                  </Button>
                ))}
              </div>
            </div>
          </Card.Body>
        </Card>

        {loading && (
          <Spinner message="Analyzing and comparing GitHub profiles. This may take a moment..." />
        )}

        {error && (
          <Alert variant="danger" className="mb-4">
            {error}
          </Alert>
        )}

        {comparison && (
          <>
            <Card className="mb-4">
              <Card.Header>
                <h3 className="mb-0">Comparison Summary</h3>
              </Card.Header>
              <Card.Body>
                <p className="lead">{comparison.comparativeSummary}</p>

                <div className="mt-4 p-3 bg-light rounded border">
                  <h4>Team Fit Analysis</h4>
                  <p>{comparison.teamFitAnalysis}</p>
                </div>
              </Card.Body>
            </Card>

            <Tabs defaultActiveKey="skills" className="mb-4" fill>
              <Tab eventKey="skills" title="Skill Comparison">
                <Card>
                  <Card.Body>
                    <h4>Technical Skill Comparison</h4>
                    {renderSkillComparison()}
                  </Card.Body>
                </Card>
              </Tab>
              <Tab eventKey="projects" title="Project Scale">
                <Card>
                  <Card.Body>
                    <h4>Project Scale Distribution</h4>
                    {renderProjectScaleComparison()}
                  </Card.Body>
                </Card>
              </Tab>
              <Tab eventKey="style" title="Collaboration Style">
                <Card>
                  <Card.Body>
                    <h4>Collaboration Style Comparison</h4>
                    <ListGroup variant="flush">
                      {Object.entries(comparison.collaborationStyle).map(
                        ([developer, style], index) => (
                          <ListGroup.Item key={index}>
                            <h5>{developer}</h5>
                            <p>{style}</p>
                          </ListGroup.Item>
                        )
                      )}
                    </ListGroup>
                  </Card.Body>
                </Card>
              </Tab>
              <Tab eventKey="learning" title="Learning Trajectory">
                <Card>
                  <Card.Body>
                    <h4>Learning Trajectory & Adaptability</h4>
                    <div className="mb-4">
                      {Object.entries(comparison.learningTrajectory).map(
                        ([developer, score], index) => (
                          <div key={index} className="mb-3">
                            <div className="d-flex justify-content-between mb-1">
                              <div>{developer}</div>
                              <div>{score.toFixed(1)}/100</div>
                            </div>
                            <ProgressBar
                              now={score}
                              variant={getColorForIndex(index)}
                              style={{ height: "10px" }}
                            />
                          </div>
                        )
                      )}
                    </div>
                  </Card.Body>
                </Card>
              </Tab>
              <Tab eventKey="complementarity" title="Complementarity">
                <Card>
                  <Card.Body>
                    <h4>Complementary Skill Areas</h4>
                    {renderComplementarityChart()}
                  </Card.Body>
                </Card>
              </Tab>
            </Tabs>
          </>
        )}
      </Container>
    </div>
  );
}

export default DeveloperComparison;
