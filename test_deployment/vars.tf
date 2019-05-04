variable "server_location" {
  default = "nbg1"
}

variable "server_type" {
  default = "cx11"
}

variable "server_image" {
  default = "debian-9"
}

variable "domain" {
}

variable "primary_pubkey" {
}

variable "admin_user" {
}

resource "random_id" "subdomain_id" {
  byte_length = 4
}

resource "random_id" "auth_key" {
  byte_length = 16
}

locals {
  subdomain = "test-${random_id.subdomain_id.hex}"
  hostname  = "${local.subdomain}.${var.domain}"
  auth_key  = "${random_id.auth_key.hex}"
}

output "hostname" {
  value = "${local.hostname}"
}
