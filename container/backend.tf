terraform {
  backend "local" {
    path = "/infra/state/terraform.tfstate"
  }
}
