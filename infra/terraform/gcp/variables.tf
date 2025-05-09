variable "service_name" {
  description = "The name of the Cloud Run service."
  type        = string
}

variable "image" {
  description = "The container image to deploy."
  type        = string
}

variable "region" {
  description = "The region where the Cloud Run service will be deployed."
  type        = string
}

variable "project_id" {
  description = "The Google Cloud project ID."
  type        = string
}