output "cloud_run_url" {
  value = google_cloud_run_service.my_service.status[0].url
}

output "cloud_run_service_name" {
  value = google_cloud_run_service.my_service.name
}