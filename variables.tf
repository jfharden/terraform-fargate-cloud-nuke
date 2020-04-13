variable "cloudwatch_build_schedule" {
  type        = "string"
  description = "Cloudwatch Schedule expression to specify when to run a build of the container. If blank then no build schdule will be created."
  default     = ""
}

variable "cloudwatch_run_schedule" {
  type        = "string"
  description = "Cloudwatch Schedule expression to specify when to run the cloud nuke task."
}

variable "cloud_nuke_version" {
  type        = "string"
  description = "Version number of cloud nuke to install."
  default     = "v0.1.17"
}

variable "name" {
  type        = string
  description = "Name to apply to all the resources."
}
