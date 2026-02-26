import { useEffect, useState } from "react"
import {
  fetchServices,
  fetchServiceDashboard,
  deployService,
} from "../api/services"
import DeployButton from "../components/DeployButton"
import ServiceCard from "../components/ServiceCard"

const DEFAULT_ENVS = ["dev", "test", "prod"]

export default function ServiceDashboard() {
  const [services, setServices] = useState([])
  const [dashboards, setDashboards] = useState({})
  const [deploying, setDeploying] = useState({})

  useEffect(() => {
    init()
  }, [])

  async function init() {
    const serviceList = await fetchServices()
    setServices(serviceList)

    const dashboardsMap = {}

    for (const svc of serviceList) {
      const dashboard = await fetchServiceDashboard(
        svc.serviceName
      )
      dashboardsMap[svc.serviceName] = dashboard // can be null
    }

    setDashboards(dashboardsMap)
  }

  const handleDeploy = async (serviceName, env) => {
    const key = `${serviceName}-${env}`
    setDeploying((p) => ({ ...p, [key]: true }))

    try {
      await deployService(serviceName, env)
      alert(`Deployment triggered: ${serviceName} → ${env}`)
    } finally {
      setDeploying((p) => ({ ...p, [key]: false }))
      await init() // refresh
    }
  }

  return (
    <div>
      <h2>Service Dashboard</h2>

      {services.map((svc) => {
        const dashboard = dashboards[svc.serviceName]
        const envs =
          dashboard?.environments
            ? Object.keys(dashboard.environments)
            : DEFAULT_ENVS

        return (
          <div
            key={svc.serviceName}
            style={{
              border: "1px solid #ccc",
              padding: 16,
              marginBottom: 16,
              borderRadius: 8,
            }}
          >
            {/* SERVICE NAME */}
            <h3>{svc.serviceName}</h3>

            {/* DEPLOY BUTTONS (ALWAYS SHOWN) */}
            <div style={{ marginBottom: 12 }}>
              {envs.map((env) => (
                <DeployButton
                  key={env}
                  env={env}
                  status={
                    dashboard?.environments?.[env]?.status ||
                    "not_deployed"
                  }
                  loading={
                    deploying[
                      `${svc.serviceName}-${env}`
                    ]
                  }
                  onDeploy={() =>
                    handleDeploy(svc.serviceName, env)
                  }
                />
              ))}
            </div>

            {/* DASHBOARD DETAILS */}
            {dashboard ? (
              <ServiceCard
                serviceName={svc.serviceName}
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
      })}
    </div>
  )
}