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
  ProgressBar,
  Tabs,
  Tab,
} from "react-bootstrap";
import ReactMarkdown from "react-markdown";

// Components
import Spinner from "../components/Spinner";

function DeveloperAnalysis() {
  const [username, setUsername] = useState("");
  const [depth, setDepth] = useState("medium");
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const API_BASE_URL =
    process.env.REACT_APP_API_BASE_URL || "http://localhost:5000/api";

  const depthOptions = [
    { value: "light", label: "Light (Faster)" },
    { value: "medium", label: "Medium (Balanced)" },
    { value: "full", label: "Full (Comprehensive)" },
  ];

  const popularDevs = [
    "torvalds",
    "gaearon",
    "tj",
    "yyx990803",
    "sindresorhus",
  ];

  const handleSuggestionClick = (suggestion) => {
    setUsername(suggestion);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!username.trim()) {
      setError("Please enter a GitHub username");
      return;
    }

    try {
      setLoading(true);
      setError(null);
      setProfile(null);

      const response = await axios.post(`${API_BASE_URL}/developers/analyze`, {
        username,
        depth,
      });

      if (response.data && response.data.success) {
        setProfile(response.data.data.profile);
      } else {
        throw new Error("Failed to analyze developer profile");
      }
    } catch (err) {
      setError(
        err.response?.data?.error ||
          err.message ||
          "An error occurred during analysis"
      );
    } finally {
      setLoading(false);
    }
  };

  const renderLanguageChart = () => {
    if (
      !profile ||
      !profile.languageAnalysis ||
      profile.languageAnalysis.length === 0
    ) {
      return <p>No language data available</p>;
    }

    return (
      <div>
        {profile.languageAnalysis.map((lang, index) => (
          <div key={index} className="mb-3">
            <div className="d-flex justify-content-between mb-1">
              <div>
                <strong>{lang.name}</strong>{" "}
                <Badge bg="secondary">{lang.proficiency}</Badge>
              </div>
              <div>{lang.percentage.toFixed(1)}%</div>
            </div>
            <ProgressBar
              now={lang.percentage}
              variant={getLanguageColor(index)}
              style={{ height: "10px" }}
            />
            <small className="text-muted">
              ~{lang.experience.toFixed(1)} years of experience
            </small>
          </div>
        ))}
      </div>
    );
  };

  const getLanguageColor = (index) => {
    const colors = ["primary", "success", "info", "warning", "danger"];
    return colors[index % colors.length];
  };

  const renderProjectComplexity = () => {
    if (
      !profile ||
      !profile.projectAnalysis ||
      profile.projectAnalysis.length === 0
    ) {
      return <p>No project data available</p>;
    }

    // Count projects by complexity
    const complexityCounts = {
      toy: 0,
      small: 0,
      medium: 0,
      large: 0,
      enterprise: 0,
    };

    profile.projectAnalysis.forEach((project) => {
      if (
        complexityCounts.hasOwnProperty(project.complexityRank.toLowerCase())
      ) {
        complexityCounts[project.complexityRank.toLowerCase()]++;
      }
    });

    return (
      <div>
        <div className="mb-4">
          <div className="d-flex justify-content-between mb-1">
            <div>Toy/Learning Projects</div>
            <div>{complexityCounts.toy}</div>
          </div>
          <ProgressBar
            now={complexityCounts.toy * 10}
            variant="info"
            style={{ height: "10px" }}
          />
        </div>
        <div className="mb-4">
          <div className="d-flex justify-content-between mb-1">
            <div>Small Projects</div>
            <div>{complexityCounts.small}</div>
          </div>
          <ProgressBar
            now={complexityCounts.small * 10}
            variant="success"
            style={{ height: "10px" }}
          />
        </div>
        <div className="mb-4">
          <div className="d-flex justify-content-between mb-1">
            <div>Medium Projects</div>
            <div>{complexityCounts.medium}</div>
          </div>
          <ProgressBar
            now={complexityCounts.medium * 10}
            variant="primary"
            style={{ height: "10px" }}
          />
        </div>
        <div className="mb-4">
          <div className="d-flex justify-content-between mb-1">
            <div>Large Projects</div>
            <div>{complexityCounts.large}</div>
          </div>
          <ProgressBar
            now={complexityCounts.large * 10}
            variant="warning"
            style={{ height: "10px" }}
          />
        </div>
        <div className="mb-4">
          <div className="d-flex justify-content-between mb-1">
            <div>Enterprise Projects</div>
            <div>{complexityCounts.enterprise}</div>
          </div>
          <ProgressBar
            now={complexityCounts.enterprise * 10}
            variant="danger"
            style={{ height: "10px" }}
          />
        </div>
      </div>
    );
  };

  const renderSkillAssessment = () => {
    if (!profile || !profile.skillAssessment) {
      return <p>No skill assessment data available</p>;
    }

    const skills = [
      { name: "Code Quality", value: profile.skillAssessment.codeQuality },
      { name: "Documentation", value: profile.skillAssessment.documentation },
      { name: "Testing", value: profile.skillAssessment.testing },
      { name: "Performance", value: profile.skillAssessment.performance },
      { name: "Security", value: profile.skillAssessment.security },
      { name: "Collaboration", value: profile.skillAssessment.collaboration },
    ];

    return (
      <div>
        {skills.map((skill, index) => (
          <div key={index} className="mb-3">
            <div className="d-flex justify-content-between mb-1">
              <div>{skill.name}</div>
              <div>{skill.value.toFixed(1)}/100</div>
            </div>
            <ProgressBar
              now={skill.value}
              variant={getSkillColor(skill.value)}
              style={{ height: "10px" }}
            />
          </div>
        ))}
      </div>
    );
  };

  const getSkillColor = (value) => {
    if (value >= 80) return "success";
    if (value >= 60) return "primary";
    if (value >= 40) return "info";
    if (value >= 20) return "warning";
    return "danger";
  };

  const renderMonthlyActivity = () => {
    if (
      !profile ||
      !profile.contributionPatterns ||
      !profile.contributionPatterns.monthlyActivity
    ) {
      return <p>No activity data available</p>;
    }

    const months = [
      "Jan",
      "Feb",
      "Mar",
      "Apr",
      "May",
      "Jun",
      "Jul",
      "Aug",
      "Sep",
      "Oct",
      "Nov",
      "Dec",
    ];
    const currentMonth = new Date().getMonth();
    const monthlyData = profile.contributionPatterns.monthlyActivity;

    // Get max value for scaling
    const maxValue = Math.max(...monthlyData);

    return (
      <div className="mt-3">
        <div className="d-flex justify-content-between">
          {monthlyData.map((count, index) => {
            // Calculate the month label based on current month
            const monthIndex = (12 + currentMonth - index) % 12;
            return (
              <div
                key={index}
                className="d-flex flex-column align-items-center"
                style={{ flex: 1 }}
              >
                <div className="text-center" style={{ fontSize: "0.8rem" }}>
                  {months[monthIndex]}
                </div>
                <div
                  className="bg-primary mt-1"
                  style={{
                    width: "80%",
                    height: `${(count / maxValue) * 100}px`,
                    minHeight: "5px",
                    borderRadius: "3px",
                  }}
                ></div>
                <div
                  className="text-center mt-1"
                  style={{ fontSize: "0.7rem" }}
                >
                  {count}
                </div>
              </div>
            );
          })}
        </div>
      </div>
    );
  };

  return (
    <div className="developer-analysis">
      <Container className="py-4">
        <div className="text-center mb-4">
          <h1>Developer Analysis</h1>
          <p className="lead">
            Analyze GitHub developers to understand their technical skills,
            project complexity, and coding patterns.
          </p>
        </div>

        <Card className="mb-4">
          <Card.Body>
            <Form onSubmit={handleSubmit}>
              <Form.Group className="mb-3">
                <Form.Label>GitHub Username</Form.Label>
                <Form.Control
                  type="text"
                  placeholder="Enter a GitHub username"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  required
                />
              </Form.Group>

              <Form.Group className="mb-3">
                <Form.Label>Analysis Depth</Form.Label>
                <Form.Select
                  value={depth}
                  onChange={(e) => setDepth(e.target.value)}
                >
                  {depthOptions.map((option) => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </Form.Select>
                <Form.Text className="text-muted">
                  Higher depth provides more detailed analysis but takes longer
                </Form.Text>
              </Form.Group>

              <Button
                type="submit"
                variant="primary"
                className="w-100"
                disabled={loading || !username.trim()}
              >
                {loading ? "Analyzing..." : "Analyze Developer"}
              </Button>
            </Form>

            <div className="mt-4">
              <h5>Popular Developers:</h5>
              <div className="d-flex flex-wrap gap-2">
                {popularDevs.map((suggestion, index) => (
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
          <Spinner message="Analyzing GitHub profile. This may take a moment..." />
        )}

        {error && (
          <Alert variant="danger" className="mb-4">
            {error}
          </Alert>
        )}

        {profile && (
          <>
            <Card className="mb-4">
              <Card.Header>
                <h3 className="mb-0">Developer Profile: {profile.username}</h3>
              </Card.Header>
              <Card.Body>
                <h4>Executive Summary</h4>
                <p className="lead">{profile.executiveSummary}</p>

                <Row className="mt-4">
                  <Col md={6} className="mb-3">
                    <h5>Technical Strengths</h5>
                    <ul className="list-group">
                      {profile.technicalStrengths.map((strength, index) => (
                        <li
                          key={index}
                          className="list-group-item list-group-item-success"
                        >
                          {strength}
                        </li>
                      ))}
                    </ul>
                  </Col>
                  <Col md={6} className="mb-3">
                    <h5>Growth Areas</h5>
                    <ul className="list-group">
                      {profile.growthAreas.map((area, index) => (
                        <li
                          key={index}
                          className="list-group-item list-group-item-warning"
                        >
                          {area}
                        </li>
                      ))}
                    </ul>
                  </Col>
                </Row>
              </Card.Body>
            </Card>

            <Tabs defaultActiveKey="languages" className="mb-4" fill>
              <Tab eventKey="languages" title="Languages">
                <Card>
                  <Card.Body>
                    <h4>Language Proficiency</h4>
                    {renderLanguageChart()}
                  </Card.Body>
                </Card>
              </Tab>
              <Tab eventKey="projects" title="Projects">
                <Card>
                  <Card.Body>
                    <h4>Project Complexity</h4>
                    {renderProjectComplexity()}
                  </Card.Body>
                </Card>
              </Tab>
              <Tab eventKey="skills" title="Skills">
                <Card>
                  <Card.Body>
                    <h4>Technical Skill Assessment</h4>
                    {renderSkillAssessment()}
                  </Card.Body>
                </Card>
              </Tab>
              <Tab eventKey="activity" title="Activity">
                <Card>
                  <Card.Body>
                    <h4>Contribution Activity</h4>
                    <div className="mb-3">
                      <Row>
                        <Col md={6}>
                          <div className="d-flex flex-column border rounded p-3 mb-3">
                            <div className="text-center mb-2">
                              Total Commits
                            </div>
                            <div className="text-center h3">
                              {profile.contributionPatterns.totalCommits}
                            </div>
                          </div>
                        </Col>
                        <Col md={6}>
                          <div className="d-flex flex-column border rounded p-3 mb-3">
                            <div className="text-center mb-2">
                              Monthly Average
                            </div>
                            <div className="text-center h3">
                              {profile.contributionPatterns.averageCommitsMonth.toFixed(
                                1
                              )}
                            </div>
                          </div>
                        </Col>
                      </Row>
                      <Row>
                        <Col md={6}>
                          <div className="d-flex flex-column border rounded p-3 mb-3">
                            <div className="text-center mb-2">
                              Consistency Score
                            </div>
                            <div className="text-center h3">
                              {profile.contributionPatterns.consistencyScore.toFixed(
                                1
                              )}
                              /100
                            </div>
                          </div>
                        </Col>
                        <Col md={6}>
                          <div className="d-flex flex-column border rounded p-3 mb-3">
                            <div className="text-center mb-2">
                              PR Acceptance Rate
                            </div>
                            <div className="text-center h3">
                              {profile.contributionPatterns.prAcceptanceRate.toFixed(
                                1
                              )}
                              %
                            </div>
                          </div>
                        </Col>
                      </Row>
                    </div>
                    <h5>Monthly Activity (Last 12 Months)</h5>
                    {renderMonthlyActivity()}
                  </Card.Body>
                </Card>
              </Tab>
              <Tab eventKey="specialization" title="Specialization">
                <Card>
                  <Card.Body>
                    <h4>Technical Specialization</h4>
                    <div className="mb-4">
                      <h5>Primary Domains</h5>
                      <div className="d-flex flex-wrap gap-2 mb-3">
                        {profile.specializations.primaryDomains.map(
                          (domain, index) => (
                            <Badge key={index} bg="primary" className="p-2">
                              {domain}
                            </Badge>
                          )
                        )}
                      </div>

                      <h5>Secondary Domains</h5>
                      <div className="d-flex flex-wrap gap-2">
                        {profile.specializations.secondaryDomains.map(
                          (domain, index) => (
                            <Badge key={index} bg="secondary" className="p-2">
                              {domain}
                            </Badge>
                          )
                        )}
                      </div>
                    </div>

                    <h5>Framework Expertise</h5>
                    <div className="mt-3">
                      {Object.entries(
                        profile.specializations.frameworkExpertise
                      ).map(([framework, value], index) => (
                        <div key={index} className="mb-3">
                          <div className="d-flex justify-content-between mb-1">
                            <div>{framework}</div>
                            <div>{value.toFixed(1)}/100</div>
                          </div>
                          <ProgressBar
                            now={value}
                            variant="info"
                            style={{ height: "10px" }}
                          />
                        </div>
                      ))}
                    </div>
                  </Card.Body>
                </Card>
              </Tab>
              <Tab eventKey="learning" title="Learning">
                <Card>
                  <Card.Body>
                    <h4>Learning Velocity</h4>
                    <div className="d-flex flex-column border rounded p-3 mb-4">
                      <div className="text-center mb-2">Adaptability Score</div>
                      <div className="text-center h2">
                        {profile.learningVelocity.adaptabilityScore.toFixed(1)}
                        /100
                      </div>
                      <div className="text-center">
                        <Badge bg="info" className="p-2">
                          Growth Trajectory:{" "}
                          {profile.learningVelocity.growthTrajectory}
                        </Badge>
                      </div>
                    </div>

                    <h5>Technology Adoption Timeline</h5>
                    <div className="mt-3">
                      <ul className="list-group">
                        {Object.entries(
                          profile.learningVelocity.technologyAdoption
                        ).map(([tech, time], index) => (
                          <li
                            key={index}
                            className="list-group-item d-flex justify-content-between align-items-center"
                          >
                            <span>{tech}</span>
                            <Badge bg="primary" pill>
                              {time}
                            </Badge>
                          </li>
                        ))}
                      </ul>
                    </div>
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

export default DeveloperAnalysis;
