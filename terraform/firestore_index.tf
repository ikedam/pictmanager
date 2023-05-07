resource "google_firestore_index" "Image-TagList-PublishTime" {
  collection = "Image"

  fields {
    field_path   = "TagList"
    array_config = "CONTAINS"
  }
  fields {
    field_path = "PublishTime"
    order      = "DESCENDING"
  }
}

resource "google_firestore_index" "Image-ItemMask-PublishTime" {
  collection = "Image"

  fields {
    field_path = "ItemMask"
    order      = "ASCENDING"
  }
  fields {
    field_path = "PublishTime"
    order      = "DESCENDING"
  }
}
