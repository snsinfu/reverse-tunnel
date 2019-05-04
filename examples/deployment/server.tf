resource "hcloud_server" "test" {
  name        = "${local.hostname}"
  server_type = "${var.server_type}"
  image       = "${var.server_image}"
  location    = "${var.server_location}"
  ssh_keys    = ["${hcloud_ssh_key.primary.id}"]
  user_data   = "${data.template_file.cloudinit.rendered}"
}

resource "hcloud_ssh_key" "primary" {
  name       = "${local.subdomain} primary"
  public_key = "${var.primary_pubkey}"
}

data "template_file" "cloudinit" {
  template = "${file("files/cloudinit.yml.in")}"
  vars {
    admin_user = "${var.admin_user}"
  }
}

data "template_file" "ssh_config" {
  template = "${file("files/ssh_config.in")}"
  vars {
    host_address = "${hcloud_server.test.ipv4_address}"
    admin_user   = "${var.admin_user}"
  }
}

data "template_file" "hosts" {
  template = "${file("files/hosts.in")}"
  vars {
    host_address = "${hcloud_server.test.ipv4_address}"
    admin_user   = "${var.admin_user}"
  }
}

output "ip" {
  value = "${hcloud_server.test.ipv4_address}"
}

output "ssh_config" {
  value     = "${data.template_file.ssh_config.rendered}"
  sensitive = true
}

output "hosts" {
  value     = "${data.template_file.hosts.rendered}"
  sensitive = true
}
