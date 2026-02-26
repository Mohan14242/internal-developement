import { useEffect, useState } from "react"
import { useParams, Link } from "react-router-dom"
import {
  fetchServiceDashboard,
  deployService,
} from "../api/services"
import DeployButton from "../components/DeployButton"
import ServiceCard from "../components/ServiceCard"

const DEFAULT_ENVS = ["dev", "test", "prod"]
const POLL_INTERVAL_MS = 5000 // 5 seconds

export default function ServiceDashboard() {
  const { serviceName } = useParams()

  const [dashboard, setDashboard] = useState(null)
  const [deploying, setDeploying] = useState({})
  const [loading, setLoading] = useState(true)

  // -------- Load dashboard once + polling ----------
  useEffect(() => {
    let isMounted = true

    async function load() {
      try {
        const data = await fetchServiceDashboard(serviceName)
        if (isMounted) setDashboard(data)
      } catch {
        if (isMounted) setDashboard(null)
      } finally {
        if (isMounted) setLoading(false)
      }
    }

    load()

    const interval = setInterval(load, POLL_INTERVAL_MS)

    return () => {
      isMounted = false
      clearInterval(interval)
    }
  }, [serviceName])

  // -------- Deploy handler ----------
  const handleDeploy = async (env) => {
    setDeploying((p) => ({ ...p, [env]: true }))

    try {
      await deployService(serviceName, env)
      alert(`Deployment triggered: ${serviceName} → ${env}`)
    } catch {
      alert("Failed to trigger deployment")
    } finally {
      setDeploying((p) => ({ ...p, [env]: false }))
    }
  }

  if (loading) {
    return <p>Loading {serviceName} dashboard...</p>
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
        <p style={{ color: "#777", fontSize: 12 }}>
          Auto-refresh every {POLL_INTERVAL_MS / 1000}s
        </p>
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

      {/* SERVICE DETAILS + ROLLBACK */}
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