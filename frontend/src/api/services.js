export async function fetchServices() {
  const res = await fetch("/api/services")
  if (!res.ok) throw new Error("Failed to fetch services")
  return res.json()
}


// src/api/services.js
export async function fetchServiceDashboard(serviceName) {
  const res = await fetch(`/api/servicesdashboard/${serviceName}/dashboard`, {
    headers: { Accept: "application/json" },
  })
  if (!res.ok) throw new Error("Failed to fetch service dashboard")
  return res.json()
}

export async function deployService(serviceName, environment) {
  const res = await fetch(`/api/deploy-services/${serviceName}/deploy`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ environment }),
  })
  if (!res.ok) throw new Error("Deployment failed")
  return res.json()
}


