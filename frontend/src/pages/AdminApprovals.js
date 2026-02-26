import { useEffect, useState } from "react"
import {
  fetchProdApprovals,
  approveDeployment,
  rejectDeployment,
} from "../api/approvals"

export default function AdminApprovals() {
  const [approvals, setApprovals] = useState([])
  const [loading, setLoading] = useState(true)
  const [actionLoading, setActionLoading] = useState({})
  const [error, setError] = useState("")

  useEffect(() => {
    loadApprovals()
  }, [])

  async function loadApprovals() {
    try {
      setLoading(true)
      setError("")

      const data = await fetchProdApprovals()
      setApprovals(data)

      console.log("[UI] Loaded approvals:", data)
    } catch (err) {
      console.error("[UI] Failed to load approvals", err)
      setError("Failed to load approvals")
    } finally {
      setLoading(false)
    }
  }

  async function handleApprove(id) {
    setActionLoading((p) => ({ ...p, [id]: true }))

    try {
      await approveDeployment(id)
      console.log("[UI] Approved deployment:", id)
      await loadApprovals()
    } catch (err) {
      console.error("[UI] Approval failed", err)
      alert("Approval failed")
    } finally {
      setActionLoading((p) => ({ ...p, [id]: false }))
    }
  }

  async function handleReject(id) {
    setActionLoading((p) => ({ ...p, [id]: true }))

    try {
      await rejectDeployment(id)
      console.log("[UI] Rejected deployment:", id)
      await loadApprovals()
    } catch (err) {
      console.error("[UI] Rejection failed", err)
      alert("Rejection failed")
    } finally {
      setActionLoading((p) => ({ ...p, [id]: false }))
    }
  }

  if (loading) {
    return <p>Loading prod approvals...</p>
  }

  if (error) {
    return <p style={{ color: "red" }}>{error}</p>
  }

  return (
    <div>
      <h2>🔐 Production Deployment Approvals</h2>

      {approvals.length === 0 && (
        <p style={{ color: "#666" }}>
          No pending production approvals 🎉
        </p>
      )}

      {approvals.map((a) => (
        <div
          key={a.id}
          style={{
            border: "1px solid #ccc",
            padding: 16,
            marginBottom: 12,
            borderRadius: 6,
          }}
        >
          <p>
            <strong>Service:</strong> {a.serviceName}
            <br />
            <strong>Version:</strong> {a.version}
            <br />
            <strong>Requested By:</strong>{" "}
            {a.requestedBy || "unknown"}
            <br />
            <strong>Requested At:</strong>{" "}
            {a.createdAt
              ? new Date(a.createdAt).toLocaleString()
              : "N/A"}
          </p>

          <button
            onClick={() => handleApprove(a.id)}
            disabled={actionLoading[a.id]}
            style={{
              marginRight: 8,
              background: "green",
              color: "white",
            }}
          >
            {actionLoading[a.id] ? "Processing..." : "Approve"}
          </button>

          <button
            onClick={() => handleReject(a.id)}
            disabled={actionLoading[a.id]}
            style={{
              background: "red",
              color: "white",
            }}
          >
            {actionLoading[a.id] ? "Processing..." : "Reject"}
          </button>
        </div>
      ))}
    </div>
  )
}