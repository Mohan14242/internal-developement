// src/components/ServiceCard.jsx
import { useState } from "react"
import { fetchArtifactsByEnv, rollbackService } from "../api/serviceApi"

export default function ServiceCard({ serviceName, dashboard }) {
  const [selectedEnv, setSelectedEnv] = useState("")
  const [artifacts, setArtifacts] = useState([])
  const [selectedVersion, setSelectedVersion] = useState("")
  const [loading, setLoading] = useState(false)

  const environments = Object.keys(dashboard.environments || {})

  const handleEnvSelect = async (env) => {
    setSelectedEnv(env)
    setSelectedVersion("")
    setLoading(true)
    try {
      const data = await fetchArtifactsByEnv(serviceName, env)
      setArtifacts(data.artifacts || [])
    } finally {
      setLoading(false)
    }
  }

  const handleRollback = async () => {
    await rollbackService(serviceName, {
      environment: selectedEnv,
      version: selectedVersion,
    })
    alert("Rollback triggered")
  }

  return (
    <div style={{ marginTop: 16 }}>
      <h3>{serviceName}</h3>

      <p>
        <strong>Owner:</strong> {dashboard.ownerTeam} <br />
        <strong>Runtime:</strong> {dashboard.runtime}
      </p>

      <strong>Rollback</strong>
      <div style={{ marginTop: 6 }}>
        {environments.map((env) => (
          <button
            key={env}
            onClick={() => handleEnvSelect(env)}
            style={{ marginRight: 8 }}
          >
            {env.toUpperCase()}
          </button>
        ))}
      </div>

      {selectedEnv && (
        <>
          <p style={{ marginTop: 8 }}>
            <strong>Current Version:</strong>{" "}
            {dashboard.environments[selectedEnv]?.currentVersion || "N/A"}
          </p>

          <select
            value={selectedVersion}
            onChange={(e) => setSelectedVersion(e.target.value)}
          >
            <option value="">Select version</option>
            {artifacts.map((a) => (
              <option key={a.version} value={a.version}>
                {a.version} —{" "}
                {a.deployedAt
                  ? new Date(a.deployedAt).toLocaleString()
                  : ""}
              </option>
            ))}
          </select>

          <br />

          <button
            disabled={!selectedVersion || loading}
            onClick={handleRollback}
            style={{ marginTop: 8 }}
          >
            Rollback
          </button>
        </>
      )}
    </div>
  )
}