// src/utils/statusColor.js
export function getStatusColor(status) {
  switch (status) {
    case "success":
      return "#2ecc71"
    case "failed":
      return "#e74c3c"
    case "deploying":
      return "#f39c12"
    default:
      return "#95a5a6"
  }
}