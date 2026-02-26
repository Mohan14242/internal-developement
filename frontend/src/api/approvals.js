export async function fetchProdApprovals() {
  console.log("[API] Fetching prod approvals")

  const res = await fetch("/api/approvals?environment=prod")

  if (!res.ok) {
    throw new Error("Failed to fetch approvals")
  }

  return res.json()
}

export async function approveDeployment(id) {
  console.log("[API] Approving deployment:", id)

  const res = await fetch(`/api/approvals/${id}/approve`, {
    method: "POST",
  })

  if (!res.ok) {
    throw new Error("Approval failed")
  }
}

export async function rejectDeployment(id) {
  console.log("[API] Rejecting deployment:", id)

  const res = await fetch(`/api/approvals/${id}/reject`, {
    method: "POST",
  })

  if (!res.ok) {
    throw new Error("Rejection failed")
  }
}