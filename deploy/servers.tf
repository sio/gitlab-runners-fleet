data "yandex_compute_image" "debian" {
  family = "debian-11"
}


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
      image_id = data.yandex_compute_image.debian.id
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
      image_id = data.yandex_compute_image.debian.id
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
  }
}
