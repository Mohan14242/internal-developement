import { useState } from "react"
import {
  fetchArtifactsByEnv,
  rollbackService,
} from "../api/serviceApi"

export default function ServiceCard({ serviceName, dashboard }) {
  const [selectedEnv, setSelectedEnv] = useState("")
  const [artifacts, setArtifacts] = useState([])
  const [selectedVersion, setSelectedVersion] = useState("")
  const [loadingArtifacts, setLoadingArtifacts] = useState(false)
  const [rollingBack, setRollingBack] = useState(false)

  const environments = Object.keys(
    dashboard.environments || {}
  )

  /* ===============================
     ENV SELECT → LOAD ARTIFACTS
     =============================== */
  const handleEnvSelect = async (env) => {
    console.info("[UI] Environment selected", {
      serviceName,
      env,
    })

    setSelectedEnv(env)
    setSelectedVersion("")
    setArtifacts([])
    setLoadingArtifacts(true)

    try {
      console.debug("[API] Fetching artifacts", {
        serviceName,
        env,
      })

      const data = await fetchArtifactsByEnv(
        serviceName,
        env
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

      console.info("[UI] Artifacts loaded", {
        serviceName,
        env,
        count: artifactsList.length,
      })
    } catch (err) {
      console.error(
        "[ERROR] Failed to load artifacts",
        err
      )
      setArtifacts([])
    } finally {
      setLoadingArtifacts(false)
    }
  }

  /* ===============================
     ROLLBACK
     =============================== */
  const handleRollback = async () => {
    if (!selectedEnv || !selectedVersion) return

    const currentVersion =
      dashboard.environments[selectedEnv]
        ?.currentVersion

    // 🚫 Prevent rollback to same version (BACKEND RULE)
    if (selectedVersion === currentVersion) {
      alert(
        "This version is already running in the selected environment"
      )
      console.warn(
        "[UI] Rollback blocked (same version)",
        {
          serviceName,
          selectedEnv,
          selectedVersion,
        }
      )
      return
    }

    console.info("[UI] Rollback initiated", {
      serviceName,
      environment: selectedEnv,
      version: selectedVersion,
    })

    setRollingBack(true)

    try {
      await rollbackService(serviceName, {
        environment: selectedEnv,
        version: selectedVersion,
      })

      alert(
        `Rollback triggered for ${serviceName} (${selectedEnv})`
      )

      console.info(
        "[API] Rollback accepted by backend",
        {
          serviceName,
          selectedEnv,
          selectedVersion,
        }
      )
    } catch (err) {
      console.error(
        "[ERROR] Rollback failed",
        err
      )
      alert("Rollback failed to start")
    } finally {
      setRollingBack(false)
    }
  }

  /* ===============================
     RENDER
     =============================== */
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

      {/* ENV BUTTONS */}
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

      {/* VERSION DROPDOWN */}
      {selectedEnv && (
        <div style={{ marginTop: 12 }}>
          <p>
            <strong>Current Version:</strong>{" "}
            {dashboard.environments[selectedEnv]
              ?.currentVersion || "N/A"}
          </p>

          <select
            value={selectedVersion}
            disabled={loadingArtifacts}
            onChange={(e) =>
              setSelectedVersion(e.target.value)
            }
            style={{ minWidth: 420 }}
          >
            <option value="">
              {loadingArtifacts
                ? "Loading versions..."
                : "Select version to rollback"}
            </option>

            {artifacts.map((a) => {
              const isCurrent =
                a.version ===
                dashboard.environments[selectedEnv]
                  ?.currentVersion

              return (
                <option
                  key={a.version}
                  value={a.version}
                  disabled={isCurrent}
                >
                  {a.version}
                  {isCurrent
                    ? " (current)"
                    : ""}{" "}
                  —{" "}
                  {a.createdAt
                    ? new Date(
                        a.createdAt
                      ).toLocaleString()
                    : ""}
                </option>
              )
            })}
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