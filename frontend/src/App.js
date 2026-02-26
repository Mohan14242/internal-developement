import { BrowserRouter, Routes, Route, Link } from "react-router-dom"
import ServicesList from "./pages/ServicesList"
import ServiceDashboard from "./pages/ServiceDashboard"
import CreateServicePage from "./pages/CreateServicePage"

export default function App() {
  return (
    <BrowserRouter>
      <header style={{ marginBottom: "20px" }}>
        <h1>🚀 Internal Developer Platform</h1>

        <nav style={{ marginBottom: "10px" }}>
          <Link to="/" style={{ marginRight: "12px" }}>
            Services
          </Link>
          <Link to="/create">Create Service</Link>
        </nav>
      </header>

      <Routes>
        {/* Service list */}
        <Route path="/" element={<ServicesList />} />

        {/* Per-service dashboard */}
        <Route
          path="/services/:serviceName"
          element={<ServiceDashboard />}
        />

        {/* Create service */}
        <Route path="/create" element={<CreateServicePage />} />
      </Routes>
    </BrowserRouter>
  )
}