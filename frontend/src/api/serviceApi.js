export const createService = async (payload) => {
  const response = await fetch("/api/create-service", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  })

  if (!response.ok) {
    throw new Error("API error")
  }

  return response.json()
}


// ✅ ADD: fetch versions by environment
export async function fetchVersionsByEnv(serviceName, env) {
  const res = await fetch(
    `/api/services/${serviceName}/artifacts?environment=${env}`
  )
  if (!res.ok) throw new Error("Failed to fetch versions")
  return res.json()
}

export async function fetchServiceEnvironments(serviceName) {
  const res = await fetch(`api/services/${serviceName}/environments`)
  if (!res.ok) throw new Error("Failed to fetch environments")
  return res.json()
}


// ✅ ADD: rollback API
export async function rollbackService(serviceName, payload) {
  const res = await fetch(`/api/services/${serviceName}/rollback`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  })
  if (!res.ok) throw new Error("Rollback failed")
}


