terraform {
  required_version = ">= 0.12"
}

module "cloud_nuke" {
  source = "../../"

  name = "cloud-nuke-daily"
}
