data "template_file" "rtun" {
  template = "${file("files/rtun.yml.in")}"
  vars {
    hostname = "${local.hostname}"
    auth_key = "${local.auth_key}"
  }
}

data "template_file" "rtun_server" {
  template = "${file("files/rtun-server.yml.in")}"
  vars {
    hostname = "${local.hostname}"
    auth_key = "${local.auth_key}"
  }
}

output "rtun" {
  value     = "${data.template_file.rtun.rendered}"
  sensitive = true
}

output "rtun_server" {
  value     = "${data.template_file.rtun_server.rendered}"
  sensitive = true
}
