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

resource "google_firestore_index" "Image-Random-LastManualTagTime" {
  collection = "Image"

  fields {
    field_path = "Random"
    order      = "ASCENDING"
  }
  fields {
    field_path = "LastManualTagTime"
    order      = "ASCENDING"
  }
}

resource "google_firestore_index" "Tag-NormalizedTo-Count" {
  collection = "Tag"

  fields {
    field_path = "NormalizedTo"
    order      = "ASCENDING"
  }
  fields {
    field_path = "Count"
    order      = "DESCENDING"
  }
}
