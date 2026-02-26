import { BrowserRouter, Routes, Route, Link } from "react-router-dom"
import ServicesList from "./pages/ServicesList"
import ServiceDashboard from "./pages/ServiceDashboard"
import CreateServicePage from "./pages/CreateServicePage"
import AdminApprovals from "./pages/AdminApprovals"

export default function App() {
  return (
    <BrowserRouter>
      <header style={{ marginBottom: "20px" }}>
        <h1>🚀 Internal Developer Platform</h1>

        <nav style={{ marginBottom: "10px" }}>
          <Link to="/" style={{ marginRight: 12 }}>
            Services
          </Link>

          <Link to="/create" style={{ marginRight: 12 }}>
            Create Service
          </Link>

          <Link to="/approvals">
            Prod Approvals
          </Link>
        </nav>
      </header>

      <Routes>
        {/* Services list */}
        <Route path="/" element={<ServicesList />} />

        {/* Per-service dashboard */}
        <Route
          path="/services/:serviceName"
          element={<ServiceDashboard />}
        />

        {/* Create service */}
        <Route path="/create" element={<CreateServicePage />} />

        {/* 🔐 Production approvals */}
        <Route
          path="/approvals"
          element={<AdminApprovals />}
        />
      </Routes>
    </BrowserRouter>
  )
}