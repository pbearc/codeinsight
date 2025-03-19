// frontend/src/components/Footer.js
import React from "react";
import { Container } from "react-bootstrap";

function Footer() {
  return (
    <footer className="bg-light py-4 mt-auto">
      <Container className="text-center">
        <p className="mb-1">CodeInsight - A Developer's Contextual Assistant</p>
        <p className="text-muted small">
          &copy; {new Date().getFullYear()} - Built with GitHub API and LLM
          Integration
        </p>
      </Container>
    </footer>
  );
}

export default Footer;
