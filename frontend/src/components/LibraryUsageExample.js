// frontend/src/components/LibraryUsageExample.js
import React, { useState } from "react";
import { Card, Button, Badge } from "react-bootstrap";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { vscDarkPlus } from "react-syntax-highlighter/dist/esm/styles/prism";

function LibraryUsageExample({ example, index, language }) {
  const [expanded, setExpanded] = useState(false);

  if (!example) return null;

  const { repository, file } = example;

  // Function to truncate code for preview
  const truncateCode = (code, maxLength = 300) => {
    if (code.length <= maxLength) return code;
    return code.substring(0, maxLength) + "...";
  };

  return (
    <Card className="mb-4 shadow-sm">
      <Card.Header>
        <div className="d-flex justify-content-between align-items-center">
          <h5 className="mb-0">Example {index + 1}</h5>
          <div>
            <a
              href={repository.url}
              target="_blank"
              rel="noopener noreferrer"
              className="me-2 text-decoration-none"
            >
              {repository.full_name}
            </a>
            {repository.stars && (
              <Badge bg="warning" text="dark">
                ⭐ {repository.stars}
              </Badge>
            )}
          </div>
        </div>
      </Card.Header>

      <Card.Body>
        <div className="d-flex justify-content-between mb-3 text-muted small">
          <span>{file.path}</span>
          <a href={file.url} target="_blank" rel="noopener noreferrer">
            View on GitHub
          </a>
        </div>

        <div className="border rounded mb-3">
          <SyntaxHighlighter
            language={language}
            style={vscDarkPlus}
            showLineNumbers={true}
            wrapLongLines={true}
          >
            {expanded ? file.content : truncateCode(file.content)}
          </SyntaxHighlighter>
        </div>

        <Button
          variant="outline-primary"
          size="sm"
          onClick={() => setExpanded(!expanded)}
        >
          {expanded ? "Show Less" : "Show More"}
        </Button>
      </Card.Body>
    </Card>
  );
}

export default LibraryUsageExample;
