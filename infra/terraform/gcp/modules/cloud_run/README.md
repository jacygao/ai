# Cloud Run Module

This module provides the necessary resources to deploy a Google Cloud Run service using Terraform.

## Usage

To use this module, include it in your Terraform configuration as follows:

```hcl
module "cloud_run" {
  source     = "./modules/cloud_run"
  service_name = "your-service-name"
  image      = "gcr.io/your-project-id/your-image:tag"
  region     = "us-central1"
}
```

## Input Variables

| Name          | Description                                   | Type   | Default       |
|---------------|-----------------------------------------------|--------|---------------|
| service_name  | The name of the Cloud Run service.           | string | n/a           |
| image         | The container image to deploy.                | string | n/a           |
| region        | The region where the Cloud Run service will be deployed. | string | "us-central1" |

## Outputs

| Name          | Description                                   |
|---------------|-----------------------------------------------|
| url           | The URL of the deployed Cloud Run service.   |