# Terraform Cloud Run Project

This project provides Terraform templates to deploy a Google Cloud Run instance. It includes a modular approach, allowing for easy configuration and reuse of the Cloud Run deployment.

## Project Structure

```
terraform-cloud-run
├── modules
│   └── cloud_run
│       ├── main.tf          # Main configuration for the Cloud Run module
│       ├── outputs.tf       # Outputs for the Cloud Run module
│       ├── variables.tf     # Input variables for the Cloud Run module
│       └── README.md        # Documentation for the Cloud Run module
├── main.tf                  # Entry point for the Terraform configuration
├── outputs.tf               # Outputs for the entire Terraform configuration
├── variables.tf             # Input variables for the entire Terraform configuration
└── README.md                # Documentation for the Terraform project
```

## Getting Started

To deploy a Cloud Run instance using this project, follow these steps:

1. **Prerequisites**
   - Ensure you have Terraform installed on your machine.
   - Set up a Google Cloud project and enable the Cloud Run API.
   - Authenticate your gcloud CLI with the necessary permissions.

2. **Configuration**
   - Modify the `variables.tf` file to set your project ID, region, and any other necessary parameters.
   - Update the `modules/cloud_run/variables.tf` file to configure the service name and container image.

3. **Deployment**
   - Navigate to the project directory.
   - Initialize Terraform: 
     ```
     terraform init
     ```
   - Plan the deployment:
     ```
     terraform plan
     ```
   - Apply the configuration:
     ```
     terraform apply
     ```

4. **Outputs**
   - After a successful deployment, the URL of the deployed Cloud Run service will be displayed in the outputs.

## Module Documentation

For detailed information on the Cloud Run module, refer to the `modules/cloud_run/README.md` file. This includes specifics on input variables and outputs related to the Cloud Run service.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.