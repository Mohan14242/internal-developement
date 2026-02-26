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
    loadServicesAndDashboards()
  }, [])

  // ✅ SAFE LOADER (NO STATE WIPE)
  async function loadServicesAndDashboards() {
    const serviceList = await fetchServices()
    setServices(serviceList)

    // 👇 update dashboards PER SERVICE
    serviceList.forEach(async (svc) => {
      try {
        const dashboard = await fetchServiceDashboard(
          svc.serviceName
        )

        setDashboards((prev) => ({
          ...prev,
          [svc.serviceName]: dashboard, // only this service
        }))
      } catch {
        setDashboards((prev) => ({
          ...prev,
          [svc.serviceName]: null,
        }))
      }
    })
  }

  const handleDeploy = async (serviceName, env) => {
    const key = `${serviceName}-${env}`

    setDeploying((prev) => ({ ...prev, [key]: true }))

    try {
      await deployService(serviceName, env)
      alert(`Deployment triggered: ${serviceName} → ${env}`)
    } finally {
      setDeploying((prev) => ({ ...prev, [key]: false }))

      // 🔁 refresh ONLY this service
      const dashboard = await fetchServiceDashboard(
        serviceName
      )

      setDashboards((prev) => ({
        ...prev,
        [serviceName]: dashboard,
      }))
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

            {/* DEPLOY BUTTONS */}
            <div style={{ marginBottom: 12 }}>
              {envs.map((env) => {
                const status =
                  dashboard &&
                  dashboard.serviceName ===
                    svc.serviceName &&
                  dashboard.environments?.[env]?.status

                return (
                  <DeployButton
                    key={env}
                    env={env}
                    status={status ?? "not_deployed"}
                    loading={
                      deploying[
                        `${svc.serviceName}-${env}`
                      ]
                    }
                    onDeploy={() =>
                      handleDeploy(svc.serviceName, env)
                    }
                  />
                )
              })}
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