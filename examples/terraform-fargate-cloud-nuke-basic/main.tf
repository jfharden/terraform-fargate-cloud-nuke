terraform {
  required_version = ">= 0.12"
}

module "cloud_nuke" {
  source = "../../"

  name = var.name
  tags = var.tags
}
