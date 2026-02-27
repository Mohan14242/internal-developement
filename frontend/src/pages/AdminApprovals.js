import { useEffect, useState } from "react"
import {
  fetchProdApprovals,
  approveDeployment,
  rejectDeployment,
} from "../api/approvals"

function statusStyle(status) {
  switch (status) {
    case "approved":
      return { color: "green", fontWeight: "bold" }
    case "rejected":
      return { color: "red", fontWeight: "bold" }
    default:
      return { color: "orange", fontWeight: "bold" }
  }
}

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
      await loadApprovals()
    } catch {
      alert("Approval failed")
    } finally {
      setActionLoading((p) => ({ ...p, [id]: false }))
    }
  }

  async function handleReject(id) {
    setActionLoading((p) => ({ ...p, [id]: true }))
    try {
      await rejectDeployment(id)
      await loadApprovals()
    } catch {
      alert("Rejection failed")
    } finally {
      setActionLoading((p) => ({ ...p, [id]: false }))
    }
  }

  if (loading) return <p>Loading prod approvals...</p>
  if (error) return <p style={{ color: "red" }}>{error}</p>

  const pending = approvals.filter((a) => a.status === "pending")
  const history = approvals.filter((a) => a.status !== "pending")

  return (
    <div>
      <h2>🔐 Production Deployment Approvals</h2>

      {/* ================= PENDING ================= */}
      <h3>🟡 Pending Approvals</h3>

      {pending.length === 0 && (
        <p style={{ color: "#666" }}>
          No pending production approvals 🎉
        </p>
      )}

      {pending.map((a) => (
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
            <strong>Requested By:</strong>{" "}
            {a.requestedBy || "unknown"}
            <br />
            <strong>Requested At:</strong>{" "}
            {new Date(a.createdAt).toLocaleString()}
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

      {/* ================= HISTORY ================= */}
      <h3 style={{ marginTop: 32 }}>📜 Approval History</h3>

      {history.length === 0 && (
        <p style={{ color: "#666" }}>
          No approval history yet
        </p>
      )}

      {history.map((a) => (
        <div
          key={a.id}
          style={{
            border: "1px solid #eee",
            padding: 14,
            marginBottom: 10,
            borderRadius: 6,
            background: "#fafafa",
          }}
        >
          <p>
            <strong>Service:</strong> {a.serviceName}
            <br />
            <strong>Status:</strong>{" "}
            <span style={statusStyle(a.status)}>
              {a.status.toUpperCase()}
            </span>
            <br />
            <strong>Requested At:</strong>{" "}
            {new Date(a.createdAt).toLocaleString()}
            {a.approvedAt && (
              <>
                <br />
                <strong>Actioned At:</strong>{" "}
                {new Date(a.approvedAt).toLocaleString()}
              </>
            )}
          </p>
        </div>
      ))}
    </div>
  )
}