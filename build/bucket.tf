// yandex_storage_bucket takes a very long time during `plan` and `update` to
// refresh its state. Therefore it's not ok to keep bucket definition in the same
// module as the rest of infra which changes a lot more often

variable "ycs3_bucket" {
  description = "Object storage bucket which will hold base VM image"
  nullable    = false
  type        = string
}

resource "yandex_storage_bucket" "images" {
  bucket   = var.ycs3_bucket
  max_size = 5 * pow(2, 30) // GB
  anonymous_access_flags {
    read        = true
    list        = false
    config_read = false
  }
  lifecycle_rule {
    id      = "Remove outdated OS images"
    enabled = true
    prefix  = ""
    expiration {
      // Older OS images may contain known security vulnerabilities.
      // CI pipeline should rotate the image long before this lifecycle rule gets
      // triggered, but in case it doesn't it's better to fail loudly and
      // explicitly with a missing image than to continue running an unsafe one
      days = 6 * 7
    }
  }
  lifecycle_rule {
    id                                     = "Clean up incomplete uploads"
    prefix                                 = ""
    enabled                                = true
    abort_incomplete_multipart_upload_days = 2
  }
  access_key = yandex_iam_service_account_static_access_key.sa-static-key.access_key
  secret_key = yandex_iam_service_account_static_access_key.sa-static-key.secret_key
}

resource "yandex_iam_service_account" "sa" {
  folder_id = local.yc_folder
  name      = "tf-service-account"
}

resource "yandex_resourcemanager_folder_iam_member" "sa-editor" {
  folder_id = local.yc_folder
  role      = "storage.editor"
  member    = "serviceAccount:${yandex_iam_service_account.sa.id}"
}

resource "yandex_iam_service_account_static_access_key" "sa-static-key" {
  service_account_id = yandex_iam_service_account.sa.id
  description        = "Terraform static key for object storage"
}

data "yandex_client_config" "client" {}

locals {
  yc_folder = data.yandex_client_config.client.folder_id
}

output "AWS_ACCESS_KEY_ID" {
  value = yandex_iam_service_account_static_access_key.sa-static-key.access_key
}

output "AWS_SECRET_ACCESS_KEY" {
  value = nonsensitive(yandex_iam_service_account_static_access_key.sa-static-key.secret_key)
}
