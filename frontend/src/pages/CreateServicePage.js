import React from "react"
import { CreateServiceFrom } from "../components/CreateServiceForm"

export const CreateServicePage = () => {
  return (
    <div className="page-container">
      <h2>Create a New Service</h2>
      <CreateServiceFrom />
    </div>
  )
}