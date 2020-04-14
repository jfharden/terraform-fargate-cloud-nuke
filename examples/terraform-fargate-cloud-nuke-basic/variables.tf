variable "name" {
  type        = string
  description = "Name to give to all resources"
}

variable "tags" {
  type        = map
  description = "Tags to apply to all resources which can be tagged"
  default     = {}
}
