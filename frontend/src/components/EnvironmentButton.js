import { getStatusColor } from "../utils/statusColor"

export default function EnvironmentButton({
  env,
  status,
  onDeploy,
  onSelect,
}) {
  return (
    <button
      onClick={() => (onSelect ? onSelect(env) : onDeploy())}
      style={{
        marginRight: "8px",
        padding: "6px 10px",
        borderRadius: "4px",
        border: "1px solid #ccc",
        backgroundColor: getStatusColor(status),
        color: "white",
        cursor: "pointer",
      }}
    >
      {onSelect ? env.toUpperCase() : `Deploy ${env}`}
    </button>
  )
}