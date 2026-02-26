import { useState } from "react"
import {
  fetchArtifactsByEnv,
  rollbackService,
} from "../api/serviceApi"

export default function ServiceCard({ serviceName, dashboard }) {
  const [selectedEnv, setSelectedEnv] = useState("")
  const [artifacts, setArtifacts] = useState([])
  const [selectedVersion, setSelectedVersion] = useState("")
  const [loadingArtifacts, setLoadingArtifacts] =
    useState(false)
  const [rollingBack, setRollingBack] =
    useState(false)

  const environments = Object.keys(
    dashboard.environments || {}
  )

  /* ================================
     ENV SELECTION + ARTIFACT FETCH
     ================================ */
  const handleEnvSelect = async (env) => {
    console.info(
      `[UI] Environment selected`,
      { serviceName, env }
    )

    setSelectedEnv(env)
    setSelectedVersion("")
    setArtifacts([])
    setLoadingArtifacts(true)

    try {
      console.debug(
        `[API] Fetching artifacts`,
        { serviceName, env }
      )

      const data = await fetchArtifactsByEnv(
        serviceName,
        env
      )

      console.debug(
        `[API] Artifacts response`,
        data
      )

      const artifactsList = Array.isArray(data)
        ? data
        : []

      artifactsList.sort(
        (a, b) =>
          new Date(b.createdAt) -
          new Date(a.createdAt)
      )

      setArtifacts(artifactsList)

      console.info(
        `[UI] Artifacts loaded`,
        {
          serviceName,
          env,
          count: artifactsList.length,
        }
      )
    } catch (err) {
      console.error(
        `[ERROR] Failed to fetch artifacts`,
        { serviceName, env, error: err }
      )
      setArtifacts([])
    } finally {
      setLoadingArtifacts(false)
    }
  }

  /* ================================
     ROLLBACK
     ================================ */
  const handleRollback = async () => {
    if (!selectedEnv) {
      console.warn(
        `[UI] Rollback blocked: no environment selected`,
        { serviceName }
      )
      alert("Please select an environment")
      return
    }

    if (!selectedVersion) {
      console.warn(
        `[UI] Rollback blocked: no version selected`,
        { serviceName, selectedEnv }
      )
      alert("Please select a version to rollback")
      return
    }

    console.info(
      `[UI] Rollback initiated`,
      {
        serviceName,
        environment: selectedEnv,
        version: selectedVersion,
      }
    )

    setRollingBack(true)

    try {
      console.debug(
        `[API] Calling rollback API`,
        {
          serviceName,
          environment: selectedEnv,
          version: selectedVersion,
        }
      )

      const response = await rollbackService(
        serviceName,
        {
          environment: selectedEnv,
          version: selectedVersion,
        }
      )

      console.info(
        `[API] Rollback accepted by backend`,
        response
      )

      alert(
        `Rollback triggered for ${serviceName} (${selectedEnv})`
      )
    } catch (err) {
      console.error(
        `[ERROR] Rollback failed to start`,
        {
          serviceName,
          environment: selectedEnv,
          version: selectedVersion,
          error: err,
        }
      )

      alert("Rollback failed to start")
    } finally {
      setRollingBack(false)
    }
  }

  /* ================================
     RENDER
     ================================ */
  return (
    <div style={{ marginTop: 20 }}>
      <h3>Service Details</h3>

      <p>
        <strong>Owner:</strong>{" "}
        {dashboard.ownerTeam || "—"}
        <br />
        <strong>Runtime:</strong>{" "}
        {dashboard.runtime || "—"}
      </p>

      <strong>Rollback</strong>

      {/* ENV SELECT */}
      <div style={{ marginTop: 8 }}>
        {environments.map((env) => (
          <button
            key={env}
            onClick={() => handleEnvSelect(env)}
            disabled={loadingArtifacts || rollingBack}
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
            onChange={(e) => {
              console.info(
                `[UI] Rollback version selected`,
                {
                  serviceName,
                  environment: selectedEnv,
                  version: e.target.value,
                }
              )
              setSelectedVersion(e.target.value)
            }}
            disabled={loadingArtifacts}
            style={{ minWidth: 420 }}
          >
            <option value="">
              {loadingArtifacts
                ? "Loading versions..."
                : "Select version to rollback"}
            </option>

            {artifacts.map((a) => (
              <option
                key={a.version}
                value={a.version}
              >
                {a.version} —{" "}
                {a.createdAt
                  ? new Date(
                      a.createdAt
                    ).toLocaleString()
                  : ""}
              </option>
            ))}
          </select>

          <br />

          <button
            onClick={handleRollback}
            disabled={
              !selectedVersion || rollingBack
            }
            style={{ marginTop: 10 }}
          >
            {rollingBack
              ? "Rolling back..."
              : "Rollback"}
          </button>
        </div>
      )}
    </div>
  )
}