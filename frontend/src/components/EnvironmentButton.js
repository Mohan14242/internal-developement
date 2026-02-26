import { getStatusColor } from "../utils/statusColor"

export default function EnvironmentButton({ env, status, onSelect }) {
  return (
    <button
      onClick={() => onSelect(env)}
      style={{
        marginRight: "8px",
        padding: "6px 10px",
        borderRadius: "4px",
        border: "1px solid #ccc",
        backgroundColor: getStatusColor(status),
        color: "white",
      }}
    >
      {env.toUpperCase()}
    </button>
  )
}