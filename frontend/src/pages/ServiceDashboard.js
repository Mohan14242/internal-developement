import { useEffect, useState } from "react"
import { useParams, Link } from "react-router-dom"
import {
  fetchServiceDashboard,
  deployService,
} from "../api/services"
import DeployButton from "../components/DeployButton"
import ServiceCard from "../components/ServiceCard"

const DEFAULT_ENVS = ["dev", "test", "prod"]

export default function ServiceDashboard() {
  const { serviceName } = useParams()

  const [dashboard, setDashboard] = useState(null)
  const [deploying, setDeploying] = useState({})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")

  useEffect(() => {
    loadDashboard()
  }, [serviceName])

  async function loadDashboard() {
    try {
      setLoading(true)
      setError("")

      const data = await fetchServiceDashboard(serviceName)
      setDashboard(data)
    } catch (err) {
      // dashboard may not exist yet (404)
      setDashboard(null)
    } finally {
      setLoading(false)
    }
  }

  const handleDeploy = async (env) => {
    setDeploying((prev) => ({ ...prev, [env]: true }))

    try {
      await deployService(serviceName, env)
      alert(`Deployment triggered: ${serviceName} → ${env}`)
    } catch {
      alert("Failed to trigger deployment")
    } finally {
      setDeploying((prev) => ({ ...prev, [env]: false }))
      await loadDashboard() // refresh only this service
    }
  }

  if (loading) {
    return <p>Loading {serviceName} dashboard...</p>
  }

  if (error) {
    return <p style={{ color: "red" }}>{error}</p>
  }

  const envs = dashboard?.environments
    ? Object.keys(dashboard.environments)
    : DEFAULT_ENVS

  return (
    <div>
      {/* HEADER */}
      <div style={{ marginBottom: 16 }}>
        <Link to="/">← Back to Services</Link>
        <h2>{serviceName} – Service Dashboard</h2>
      </div>

      {/* DEPLOY BUTTONS */}
      <div style={{ marginBottom: 12 }}>
        {envs.map((env) => (
          <DeployButton
            key={env}
            env={env}
            status={
              dashboard?.environments?.[env]?.status ||
              "not_deployed"
            }
            loading={deploying[env]}
            onDeploy={() => handleDeploy(env)}
          />
        ))}
      </div>

      {/* DASHBOARD DETAILS */}
      {dashboard ? (
        <ServiceCard
          serviceName={serviceName}
          dashboard={dashboard}
        />
      ) : (
        <p style={{ color: "#777" }}>
          Service not deployed yet.  
          Deploy to initialize dashboard.
        </p>
      )}
    </div>
  )
}