export async function fetchServices() {
  const res = await fetch("/api/services")
  if (!res.ok) throw new Error("Failed to fetch services")
  return res.json()
}

export async function deployService(serviceName, env) {
  const res = await fetch(`/api/services/${serviceName}/deploy/${env}`, {
    method: "POST",
  })
  if (!res.ok) throw new Error("Deploy failed")
}

