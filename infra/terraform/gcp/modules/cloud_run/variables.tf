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

variable "memory" {
  description = "The amount of memory allocated to the Cloud Run service."
  type        = string
  default     = "256Mi"
}

variable "timeout" {
  description = "The maximum duration for a request to complete."
  type        = number
  default     = 60
}

variable "allow_unauthenticated" {
  description = "Allow unauthenticated invocations of the Cloud Run service."
  type        = bool
  default     = true
}