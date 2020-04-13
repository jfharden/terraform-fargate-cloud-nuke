resource "aws_ecr_repository" "cloud_nuke" {
  name = var.name
  image_scanning_configuration {
    scan_on_push = true
  }
}
