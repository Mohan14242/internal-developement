import { BrowserRouter, Routes, Route, Link } from "react-router-dom"
import ServiceDashboard from "./pages/ServiceDashboard"
import CreateServicePage from "./pages/CreateServicePage"

export default function App() {
  return (
    <BrowserRouter>
      <h1>🚀 Internal Developer Platform</h1>
      <Link to="/">Dashboard</Link> | <Link to="/create">Create Service</Link>

      <Routes>
        <Route path="/" element={<ServiceDashboard />} />
        <Route path="/create" element={<CreateServicePage />} />
      </Routes>
    </BrowserRouter>
  )
}