import { getStatusColor } from "../utils/statusColor"

export default function deploybutton({ env, status, onDeploy, loading }) {
  return (
    <button
      disabled={loading}
      onClick={onDeploy}
      style={{
        marginRight: "8px",
        padding: "6px 12px",
        borderRadius: "4px",
        border: "1px solid #ccc",
        backgroundColor: loading ? "orange" : getStatusColor(status),
        color: "white",
        cursor: loading ? "not-allowed" : "pointer",
      }}
    >
      {loading ? `Deploying ${env}...` : `Deploy to ${env.toUpperCase()}`}
    </button>
  )
}