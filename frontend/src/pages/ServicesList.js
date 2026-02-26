import { useEffect, useState } from "react"
import { fetchServices } from "../api/services"
import { useNavigate } from "react-router-dom"

export default function ServicesList() {
  const [services, setServices] = useState([])
  const navigate = useNavigate()

  useEffect(() => {
    load()
  }, [])

  async function load() {
    const data = await fetchServices()
    setServices(data)
  }

  return (
    <div>
      <h2>Services</h2>

      {services.map((svc) => (
        <div
          key={svc.serviceName}
          onClick={() =>
            navigate(`/services/${svc.serviceName}`)
          }
          style={{
            border: "1px solid #ccc",
            padding: 16,
            marginBottom: 12,
            borderRadius: 8,
            cursor: "pointer",
          }}
        >
          <strong>{svc.serviceName}</strong>
        </div>
      ))}
    </div>
  )
}

