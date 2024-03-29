resource "yandex_vpc_network" "vpc" {
  count = local.one_or_none
  name  = "vpc"
}

resource "yandex_vpc_subnet" "inner" {
  count          = local.one_or_none
  name           = "inner"
  v4_cidr_blocks = var.inner_cidr
  network_id     = yandex_vpc_network.vpc[0].id
  route_table_id = yandex_vpc_route_table.behind_nat[0].id
}

resource "yandex_vpc_subnet" "outer" {
  count          = local.one_or_none
  name           = "outer"
  v4_cidr_blocks = var.outer_cidr
  network_id     = yandex_vpc_network.vpc[0].id
}

resource "yandex_vpc_route_table" "behind_nat" {
  count      = local.one_or_none
  name       = "behind_nat"
  network_id = yandex_vpc_network.vpc[0].id
  static_route {
    destination_prefix = "0.0.0.0/0"
    next_hop_address   = local.gateway_ip
  }
}
