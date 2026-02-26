import { useState } from "react"
import {
  fetchArtifactsByEnv,
  rollbackService,
} from "../api/serviceApi"

export default function ServiceCard({ serviceName, dashboard }) {
  const [selectedEnv, setSelectedEnv] = useState("")
  const [artifacts, setArtifacts] = useState([])
  const [selectedVersion, setSelectedVersion] = useState("")
  const [loading, setLoading] = useState(false)

  const environments = Object.keys(
    dashboard.environments || {}
  )

  // -------- Select environment → load versions ----------
  const handleEnvSelect = async (env) => {
    setSelectedEnv(env)
    setSelectedVersion("")
    setArtifacts([])
    setLoading(true)

    try {
      const data = await fetchArtifactsByEnv(
        serviceName,
        env
      )
      setArtifacts(data.artifacts || [])
    } catch {
      setArtifacts([])
    } finally {
      setLoading(false)
    }
  }

  // -------- Trigger rollback ----------
  const handleRollback = async () => {
    if (!selectedEnv || !selectedVersion) return

    setLoading(true)
    try {
      await rollbackService(serviceName, {
        environment: selectedEnv,
        version: selectedVersion,
      })
      alert(
        `Rollback triggered: ${serviceName} → ${selectedEnv}`
      )
    } catch {
      alert("Rollback failed to start")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ marginTop: 20 }}>
      <h3>Service Details</h3>

      <p>
        <strong>Owner:</strong> {dashboard.ownerTeam}
        <br />
        <strong>Runtime:</strong> {dashboard.runtime}
      </p>

      {/* ROLLBACK */}
      <strong>Rollback</strong>

      {/* ENV SELECT */}
      <div style={{ marginTop: 8 }}>
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

      {/* VERSION SELECT */}
      {selectedEnv && (
        <div style={{ marginTop: 12 }}>
          <p>
            <strong>Current Version:</strong>{" "}
            {dashboard.environments[selectedEnv]
              ?.currentVersion || "N/A"}
          </p>

          <select
            value={selectedVersion}
            onChange={(e) =>
              setSelectedVersion(e.target.value)
            }
            style={{ minWidth: 420 }}
          >
            <option value="">
              Select version to rollback
            </option>

            {artifacts.map((a) => (
              <option key={a.version} value={a.version}>
                {a.version}
                {a.deployedAt
                  ? ` — ${new Date(
                      a.deployedAt
                    ).toLocaleString()}`
                  : ""}
              </option>
            ))}
          </select>

          <br />

          <button
            onClick={handleRollback}
            disabled={!selectedVersion || loading}
            style={{ marginTop: 10 }}
          >
            {loading ? "Rolling back..." : "Rollback"}
          </button>
        </div>
      )}
    </div>
  )
}