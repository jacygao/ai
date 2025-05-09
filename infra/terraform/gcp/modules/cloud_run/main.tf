resource "google_cloud_run_service" "default" {
  name     = var.service_name
  location = var.region

  template {
    spec {
      containers {
        image = var.image

        resources {
          limits = {
            cpu    = "1"
            memory = "256Mi"
          }
        }
      }
    }
  }

  traffic {
    percent         = 100
    revision_name   = google_cloud_run_service.default.latest_revision
    latest_revision = true
  }
}

resource "google_cloud_run_service_iam_member" "invoker" {
  service = google_cloud_run_service.default.name
  location = var.region
  role    = "roles/run.invoker"
  member  = var.invoker_member
}