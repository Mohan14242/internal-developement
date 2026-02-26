// src/components/DeployButton.jsx
import { getStatusColor } from "../utils/statusColor"

export default function DeployButton({ env, status, loading, onDeploy }) {
  return (
    <button
      disabled={loading}
      onClick={onDeploy}
      style={{
        marginRight: 8,
        padding: "6px 12px",
        borderRadius: 4,
        border: "none",
        backgroundColor: loading ? "#f39d12fc" : getStatusColor(status),
        color: "#fff",
        cursor: loading ? "not-allowed" : "pointer",
      }}
    >
      {loading ? `Deploying ${env}` : `Deploy ${env.toUpperCase()}`}
    </button>
  )
}