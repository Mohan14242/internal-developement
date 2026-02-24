import { useEffect, useState } from "react"
import EnvironmentButton from "./EnvironmentButton"
import {
  fetchServiceEnvironments,
  fetchVersionsByEnv,
  rollbackService,
} from "../api/serviceApi"

export default function ServiceCard({ service }) {
  const [environments, setEnvironments] = useState([])
  const [selectedEnv, setSelectedEnv] = useState("")
  const [versions, setVersions] = useState([])
  const [selectedVersion, setSelectedVersion] = useState("")

  useEffect(() => {
    fetchServiceEnvironments(service.serviceName).then((data) =>
      setEnvironments(data.environments)
    )
  }, [service.serviceName])

  const handleEnvSelect = async (env) => {
    setSelectedEnv(env)
    setSelectedVersion("")
    const data = await fetchVersionsByEnv(service.serviceName, env)
    setVersions(data.versions)
  }

  const handleRollback = async () => {
    await rollbackService(service.serviceName, {
      environment: selectedEnv,
      version: selectedVersion,
    })
    alert("Rollback triggered")
  }

  return (
    <>
      <h3>{service.serviceName}</h3>

      <p>
        <strong>Owner:</strong> {service.ownerTeam} <br />
        <strong>Runtime:</strong> {service.runtime} <br />
        <strong>CICD:</strong> {service.cicdType} <br />
        <strong>Deploy Type:</strong> {service.deployType}
      </p>

      {/* ROLLBACK */}
      <div>
        <strong>Rollback Environment</strong>
        <div style={{ marginTop: "6px" }}>
          {environments.map((env) => (
            <EnvironmentButton
              key={env}
              env={env}
              status="default"
              onSelect={handleEnvSelect}
            />
          ))}
        </div>
      </div>

      {selectedEnv && (
        <>
          <p>
            <strong>Current Version:</strong>{" "}
            {service.environments?.[selectedEnv]?.currentVersion || "N/A"}
          </p>

          <select
            value={selectedVersion}
            onChange={(e) => setSelectedVersion(e.target.value)}
          >
            <option value="">Select Version</option>
            {versions.map((v) => (
              <option key={v.version} value={v.version}>
                {v.version} —{" "}
                {new Date(v.deployedAt).toLocaleString()}
              </option>
            ))}
          </select>

          <br />

          <button
            disabled={!selectedVersion}
            onClick={handleRollback}
            style={{ marginTop: "8px" }}
          >
            Rollback
          </button>
        </>
      )}
    </>
  )
}