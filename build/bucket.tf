variable "ycs3_bucket" {
  description = "Object storage bucket which will hold base VM image"
  nullable    = false
  type        = string
  sensitive   = true
}

resource "yandex_storage_bucket" "images" {
  bucket     = var.ycs3_bucket
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
