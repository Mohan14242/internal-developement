import { useState } from "react"
import yaml from "js-yaml"

const DEFAULT_YAML = `
serviceName: orders
repoName: orders-service
ownerTeam: payments
runtime: go
environment: dev
`.trim()

export const CreateServiceFrom = () => {
  const [yamlText, setYamlText] = useState(DEFAULT_YAML)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")
  const [success, setSuccess] = useState("")

  const submitYAML = async () => {
    setError("")
    setSuccess("")

    if (!yamlText.trim()) {
      setError("❌ YAML cannot be empty")
      return
    }

    // Optional YAML validation
    try {
      yaml.load(yamlText)
    } catch (e) {
      setError("❌ Invalid YAML format")
      return
    }

    setLoading(true)

    try {
      const res = await fetch("/api/create-service", {
        method: "POST",
        headers: {
          "Content-Type": "application/x-yaml",
        },
        body: yamlText,
      })

      if (!res.ok) {
        const text = await res.text()
        throw new Error(text || "API error")
      }

      const data = await res.json()
      setSuccess(`✅ Service created: ${data.repoUrl}`)
    } catch (err) {
      setError(`❌ ${err.message || "Unknown error"}`)
    } finally {
      setLoading(false)
    }
  }

  const onFileSelect = async (file) => {
    const text = await file.text()
    setYamlText(text)
    setSuccess("📄 YAML file loaded")
  }

  return (
    <div style={{ maxWidth: "900px" }}>
      <h2>Create Service (YAML)</h2>

      <textarea
        rows={18}
        style={{ width: "100%", fontFamily: "monospace" }}
        value={yamlText}
        onChange={(e) => setYamlText(e.target.value)}
      />

      <br /><br />

      <input
        type="file"
        accept=".yaml,.yml"
        onChange={(e) => {
          const file = e.target.files?.[0]
          if (file) onFileSelect(file)
        }}
      />

      <br /><br />

      <button onClick={submitYAML} disabled={loading}>
        {loading ? "Creating..." : "Create Service"}
      </button>

      {error && <p style={{ color: "red" }}>{error}</p>}
      {success && <p style={{ color: "green" }}>{success}</p>}
    </div>
  )
}