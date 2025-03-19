// frontend/src/components/Navbar.js
import React from "react";
import { Link } from "react-router-dom";
import { Navbar as BootstrapNavbar, Nav, Container } from "react-bootstrap";

function Navbar() {
  return (
    <BootstrapNavbar bg="light" expand="lg" className="shadow-sm mb-3">
      <Container>
        <BootstrapNavbar.Brand
          as={Link}
          to="/"
          className="fw-bold text-primary"
        >
          CodeInsight
        </BootstrapNavbar.Brand>
        <BootstrapNavbar.Toggle aria-controls="basic-navbar-nav" />
        <BootstrapNavbar.Collapse id="basic-navbar-nav">
          <Nav className="ms-auto">
            <Nav.Link as={Link} to="/code-analyzer">
              Code Analyzer
            </Nav.Link>
            <Nav.Link as={Link} to="/library-explorer">
              Library Explorer
            </Nav.Link>
            <Nav.Link as={Link} to="/documentation-generator">
              Documentation Generator
            </Nav.Link>
            <Nav.Link as={Link} to="/implementation-finder">
              Implementation Finder
            </Nav.Link>
          </Nav>
        </BootstrapNavbar.Collapse>
      </Container>
    </BootstrapNavbar>
  );
}

export default Navbar;
