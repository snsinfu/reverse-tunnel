resource "cloudflare_record" "test" {
  domain  = "${var.domain}"
  name    = "${local.subdomain}"
  type    = "A"
  value   = "${hcloud_server.test.ipv4_address}"
  proxied = false
}
