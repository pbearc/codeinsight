import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import "bootstrap/dist/css/bootstrap.min.css";

// Components
import Navbar from "./components/Navbar";
import Footer from "./components/Footer";

// Pages
import Home from "./pages/Home";
import CodeAnalyzer from "./pages/CodeAnalyzer";
import LibraryExplorer from "./pages/LibraryExplorer";
import DocumentationGenerator from "./pages/DocumentationGenerator";
import ImplementationFinder from "./pages/ImplementationFinder";
import ReadmeGenerator from "./pages/ReadmeGenerator";
import ProjectVisualizer from "./pages/ProjectVisualizer";
import DeveloperAnalysis from "./pages/DeveloperAnalysis";
import DeveloperComparison from "./pages/DeveloperComparison";

function App() {
  return (
    <Router>
      <div className="d-flex flex-column min-vh-100">
        <Navbar />
        <main className="flex-grow-1">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/code-analyzer" element={<CodeAnalyzer />} />
            <Route path="/library-explorer" element={<LibraryExplorer />} />
            <Route
              path="/documentation-generator"
              element={<DocumentationGenerator />}
            />
            <Route
              path="/implementation-finder"
              element={<ImplementationFinder />}
            />
            <Route path="/readme-generator" element={<ReadmeGenerator />} />
            <Route path="/project-visualizer" element={<ProjectVisualizer />} />
            <Route path="/developer-analysis" element={<DeveloperAnalysis />} />
            <Route
              path="/developer-comparison"
              element={<DeveloperComparison />}
            />
          </Routes>
        </main>
        <Footer />
      </div>
    </Router>
  );
}

export default App;
