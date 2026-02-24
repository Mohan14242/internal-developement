import { useEffect, useState } from "react"
import { fetchServices, deployService } from "../api/services"
import { fetchServiceEnvironments } from "../api/serviceApi"
import deploybutton from "../components/deploybutton"
import ServiceCard from "../components/ServiceCard"

export default function ServiceDashboard() {
  const [services, setServices] = useState([])
  const [serviceEnvs, setServiceEnvs] = useState({})
  const [deploying, setDeploying] = useState({}) // service+env loading state

  // Load services + environments
  useEffect(() => {
    loadServices()
  }, [])

  const loadServices = async () => {
    const data = await fetchServices()
    setServices(data)

    const envMap = {}
    for (const svc of data) {
      const res = await fetchServiceEnvironments(svc.serviceName)
      envMap[svc.serviceName] = res.environments
    }
    setServiceEnvs(envMap)
  }

  // Deploy handler
  const handleDeploy = async (serviceName, env) => {
    const key = `${serviceName}-${env}`
    setDeploying((prev) => ({ ...prev, [key]: true }))

    try {
      await deployService(serviceName, env)
      alert(`Deployment triggered for ${serviceName} → ${env}`)
    } catch (err) {
      alert("Deployment failed to start")
    } finally {
      setDeploying((prev) => ({ ...prev, [key]: false }))
      await loadServices() // 🔁 refresh statuses
    }
  }

  return (
    <div>
      <h2>Service Dashboard</h2>

      {services.map((service) => (
        <div
          key={service.serviceName}
          style={{
            border: "1px solid #ccc",
            padding: "16px",
            marginBottom: "16px",
            borderRadius: "8px",
          }}
        >
          {/* DEPLOY BUTTONS */}
          <div style={{ marginBottom: "12px" }}>
            <strong>Deploy:</strong>{" "}
            {serviceEnvs[service.serviceName]?.map((env) => {
              const status =
                service.environments?.[env]?.status || "default"

              const loadingKey = `${service.serviceName}-${env}`

              return (
                <deploybutton
                  key={env}
                  env={env}
                  status={status}
                  loading={deploying[loadingKey]}
                  onDeploy={() =>
                    handleDeploy(service.serviceName, env)
                  }
                />
              )
            })}
          </div>

          {/* SERVICE DETAILS + ROLLBACK */}
          <ServiceCard service={service} />
        </div>
      ))}
    </div>
  )
}