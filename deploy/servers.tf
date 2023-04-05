// Compute Cloud instance is cheaper than API gateway in our usecase,
// and it's a lot easier to configure too. VPC gateway, on the other hand,
// can't do HTTP reverse proxying (or even plain TCP port forwarding).
resource "yandex_compute_instance" "gateway" {
  count       = local.one_or_none
  name        = "gateway"
  hostname    = "gateway"
  zone        = var.yc_zone
  platform_id = var.yc_platform
  resources {
    cores         = 2
    memory        = 2
    core_fraction = 20
  }
  boot_disk {
    auto_delete = true
    initialize_params {
      image_id = yandex_compute_image.base[0].id
      size     = 10
    }
  }
  network_interface {
    subnet_id      = yandex_vpc_subnet.outer[0].id
    nat            = true
    ip_address     = local.gateway_ip
    nat_ip_address = local.external_ip
  }
  scheduling_policy {
    preemptible = false // we don't want our gateway to go down suddenly
  }
  metadata = {
    serial-port-enable = 1
    user-data = templatefile("cloud-config/gateway.yml", {
      inner_subnet = var.inner_cidr[0],
    })
  }
}


// GitLab runner hosts
resource "yandex_compute_instance" "runner" {
  for_each    = var.runners
  name        = each.key
  hostname    = each.key
  zone        = var.yc_zone
  platform_id = var.yc_platform
  resources {
    cores         = 2
    memory        = 4
    core_fraction = 20
  }
  boot_disk {
    auto_delete = true
    initialize_params {
      image_id = yandex_compute_image.base[0].id
      size     = 15
    }
  }
  network_interface {
    subnet_id = yandex_vpc_subnet.inner[0].id
  }
  scheduling_policy {
    preemptible = true
  }
  metadata = {
    serial-port-enable = 1
    user-data = templatefile("cloud-config/runner.yml", {
      gitlab_runner_tag   = var.gitlab_runner_tag,
      gitlab_runner_token = var.gitlab_runner_token,
    })
  }
}

// Creation/deletion of a base image from S3 storage takes 10-15 seconds.
// This is fast enough for me to scale down to zero when no runners are in use.
// Frequent recreation also solves the problem of keeping the base image up to
// date:
// - CI pipeline will overwrite the image in S3 bucket as often as required
// - Terraform will pick up the latest image version after each inactivity cycle
resource "yandex_compute_image" "base" {
  count       = local.one_or_none
  name        = "base"
  description = "Base image both for gateway and for runners"
  os_type     = "LINUX"
  source_url  = var.ycs3_vmimage_url
  pooled      = "false"
}
