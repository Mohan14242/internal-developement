export function getStatusColor(status) {
  switch (status) {
    case "success":
      return "green"
    case "failed":
      return "red"
    default:
      return "gray" // initial state
  }
}