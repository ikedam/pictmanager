terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
    }
  }
}



provider "google" {
  project = "ikedam"
  region  = "us-central1"
  zone    = "us-central1-a"
}

data "google_project" "project" {
}

data "google_client_config" "client" {
}

locals {
  location = split("-", data.google_client_config.client.region)[0]
}

resource "google_app_engine_application" "app" {
  location_id = substr(
    data.google_client_config.client.region,
    0,
    length(data.google_client_config.client.region) - 1,
  )
  database_type = "CLOUD_FIRESTORE"
}

resource "google_app_engine_domain_mapping" "domain_mapping" {
  domain_name = "pict.ikedam.jp"

  ssl_settings {
    ssl_management_type = "AUTOMATIC"
  }
}

resource "google_service_account" "appengine" {
  account_id = "pictmanager"
}

resource "google_project_iam_member" "project" {
  project = data.google_project.project.project_id
  role    = "roles/datastore.user"
  member  = "serviceAccount:${google_service_account.appengine.email}"
}

resource "google_storage_bucket" "picts" {
  name     = "ikadam-picts"
  location = data.google_client_config.client.region
}

resource "google_storage_bucket_iam_member" "picts-appengine" {
  bucket = google_storage_bucket.picts.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${google_service_account.appengine.email}"
}

resource "google_storage_bucket_iam_member" "picts-public" {
  bucket = google_storage_bucket.picts.name
  role   = "roles/storage.legacyObjectReader"
  member = "allUsers"
}

resource "google_storage_bucket" "registry" {
  name     = "${local.location}.artifacts.${data.google_project.project.project_id}.appspot.com"
  location = upper(local.location)

  lifecycle_rule {
    condition {
      age = 3
    }
    action {
      type = "Delete"
    }
  }
}
