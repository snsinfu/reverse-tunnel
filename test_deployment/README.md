# Hetzner Cloud + Cloudflare DNS deployment

Set environment variables:

```
export HCLOUD_TOKEN=...
export CLOUDFLARE_EMAIL=...
export CLOUDFLARE_TOKEN=...
```

Create `.vars.json`:

```
{
  "domain": "example.com",
  "primary_pubkey": "ssh-ed25519 ...",
  "admin_user": "alice"
}
```

Provision and serve:

```
make
make serve
```

Run agent:

```
./rtun
```
