variable "service_name" {
  description = "The name of the Cloud Run service."
  type        = string
  default = "hello-world"
}

variable "image" {
  description = "The container image to deploy."
  type        = string
}

variable "region" {
  description = "The region where the Cloud Run service will be deployed."
  type        = string
  default = "us"
}

variable "project_id" {
  description = "The Google Cloud project ID."
  type        = string
  default = "electric-facet-306612"
}

variable "invoker_member" {
  description = "The IAM member to grant the Cloud Run Invoker role."
  type        = string
  default = "aus.jacy@gmail.com"
}