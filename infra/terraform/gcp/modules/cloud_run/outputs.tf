output "cloud_run_url" {
  value = google_cloud_run_service.my_service.status[0].url
}

output "cloud_run_service_name" {
  value = google_cloud_run_service.my_service.name
}

output "url" {
  value = google_cloud_run_service.my_service.status[0].url
  description = "The URL of the deployed Cloud Run service."
}