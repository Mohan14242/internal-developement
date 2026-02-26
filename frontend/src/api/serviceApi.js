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




export async function fetchServiceEnvironments(serviceName) {
  const res = await fetch(`/api/services/${serviceName}/environments`)
  if (!res.ok) throw new Error("Failed to fetch environments")
  return res.json()
}


// src/api/serviceApi.js
export async function fetchArtifactsByEnv(serviceName, environment) {
  const res = await fetch(
    `/api/artifact-by-env/${serviceName}/artifacts?environment=${environment}`
  )
  if (!res.ok) throw new Error("Failed to fetch artifacts")
  return res.json()
}

export async function rollbackService(serviceName, payload) {
  const res = await fetch(
    `/api/rollback-services/${serviceName}/rollback`,
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    }
  )
  if (!res.ok) throw new Error("Rollback failed")
  return res.json()
}



