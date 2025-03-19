// frontend/src/components/AnalysisResult.js
import React from "react";
import { Card } from "react-bootstrap";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { vscDarkPlus } from "react-syntax-highlighter/dist/esm/styles/prism";
import ReactMarkdown from "react-markdown";

function AnalysisResult({ result }) {
  if (!result) return null;

  // Function to extract code blocks from markdown and highlight them
  const renderMarkdown = (content) => {
    return (
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
        {content}
      </ReactMarkdown>
    );
  };

  return (
    <Card className="shadow-sm mb-4">
      <Card.Header>
        <h4 className="mb-0">Analysis Results</h4>
      </Card.Header>
      <Card.Body>
        <div className="markdown-body">{renderMarkdown(result.analysis)}</div>
      </Card.Body>
    </Card>
  );
}

export default AnalysisResult;
