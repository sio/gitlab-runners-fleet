terraform {
  required_providers {
    yandex = {
      source = "yandex-cloud/yandex"
    }
  }
  required_version = ">= 0.13"
}

provider "yandex" {
  zone = var.yc_zone
}

variable "yc_zone" {
  type    = string
  default = "ru-central1-a"
}
