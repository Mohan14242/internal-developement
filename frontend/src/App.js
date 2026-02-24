import React from "react"
import { BrowserRouter as Router, Routes, Route, Link } from "react-router-dom"
import { CreateServicePage } from "./pages/CreateServicePage"
import ServiceDashboard from "./pages/ServiceDashboard"

function App() {
  return (
    <Router>
      <div className="app-container">
        <header style={{ marginBottom: "20px" }}>
          <h1>🚀 Internal Developer Platform</h1>

          {/* Simple Navigation */}
          <nav style={{ marginTop: "10px" }}>
            <Link to="/" style={{ marginRight: "15px" }}>
              Service Dashboard
            </Link>
            <Link to="/create">
              Create Service
            </Link>
          </nav>
        </header>

        <Routes>
          <Route path="/" element={<ServiceDashboard />} />
          <Route path="/create" element={<CreateServicePage />} />
        </Routes>
      </div>
    </Router>
  )
}

export default App