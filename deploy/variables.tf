//
// Runner instances
//
variable "runners" {
  description = "Set of instance hostnames to create/keep"
  type        = set(string)
  default     = []
  nullable    = false
}
locals {
  one_or_none = length(var.runners) > 0 ? 1 : 0
}


//
// Networks
//
variable "inner_cidr" {
  type    = list(string)
  default = ["10.10.100.0/24"]
}
variable "outer_cidr" {
  type    = list(string)
  default = ["10.10.10.0/24"]
}
locals {
  gateway_ip = cidrhost(var.outer_cidr[0], 10) // .1 is reserved by YCloud for itself
}


//
// Yandex Cloud
//
variable "yc_zone" {
  type    = string
  default = "ru-central1-a"
}
variable "yc_platform" {
  description = "CPU generation <https://cloud.yandex.com/en-ru/docs/compute/concepts/vm-platforms>"
  // Intel Ice Lake (Xeon Gold 6338), cheapest as of 2023-02-17:
  // https://cloud.yandex.com/en-ru/docs/compute/pricing#prices
  default = "standard-v3"
  type    = string
}
variable "ycs3_vmimage_url" {
  description = "URL to base VM image in Yandex Cloud object storage"
  nullable    = false
  type        = string
}
