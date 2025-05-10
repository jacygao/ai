resource "google_cloud_run_service" "hello-world" {
  name     = var.service_name
  location = var.region

  template {
    spec {
      containers {
        image = var.image

        resources {
          limits = {
            cpu    = "1"
            memory = "512Mi"
          }
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_service_iam_member" "invoker" {
  service = google_cloud_run_service.hello-world.name
  location = google_cloud_run_service.hello-world.location
  role    = "roles/run.invoker"
  member  = var.invoker_member
}

module "cloud_run" {
  source       = "./modules/cloud_run"
  service_name = var.service_name
  image        = var.image
  region       = var.region
  project_id   = var.project_id
  invoker_member = var.invoker_member
}