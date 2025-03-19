// frontend/src/components/CodeEditor.js
import React from "react";
import { Form } from "react-bootstrap";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { vscDarkPlus } from "react-syntax-highlighter/dist/esm/styles/prism";

function CodeEditor({ value, onChange, language }) {
  return (
    <div className="position-relative border rounded mb-3">
      <Form.Control
        as="textarea"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder="Paste your code here or type..."
        className="font-monospace"
        style={{ minHeight: "300px", resize: "vertical" }}
      />

      {value && (
        <div
          className="position-absolute top-0 start-0 w-100 h-100 overflow-auto"
          style={{ pointerEvents: "none", opacity: 0.9 }}
        >
          <SyntaxHighlighter
            language={language}
            style={vscDarkPlus}
            showLineNumbers={true}
            wrapLongLines={true}
            customStyle={{
              margin: 0,
              borderRadius: "0.375rem",
              height: "100%",
              fontSize: "14px",
            }}
          >
            {value || "// Your code preview will appear here"}
          </SyntaxHighlighter>
        </div>
      )}
    </div>
  );
}

export default CodeEditor;
