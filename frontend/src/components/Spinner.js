// frontend/src/components/Spinner.js
import React from "react";
import { Spinner as BootstrapSpinner } from "react-bootstrap";

function Spinner({ message }) {
  return (
    <div className="text-center my-5">
      <BootstrapSpinner animation="border" role="status" variant="primary">
        <span className="visually-hidden">Loading...</span>
      </BootstrapSpinner>
      <p className="mt-3 text-muted">{message || "Loading..."}</p>
    </div>
  );
}

export default Spinner;
