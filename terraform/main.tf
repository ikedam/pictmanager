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

resource "google_app_engine_application" "app" {
  location_id   = "us-central"
  database_type = "CLOUD_FIRESTORE"
}

resource "google_storage_bucket" "picts" {
  name     = "ikadam-picts"
  location = "us-central1"
}

resource "google_storage_bucket_iam_member" "picts-public" {
  bucket = google_storage_bucket.picts.name
  role   = "roles/storage.legacyObjectReader"
  member = "allUsers"
}
