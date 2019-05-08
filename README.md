Reverse tunnel TCP and UDP
==========================

[![Release][release-badge]][release-url]
[![Build Status][travis-badge]][travis-url]
[![MIT License][license-badge]][license-url]

[release-badge]: https://img.shields.io/github/release/snsinfu/reverse-tunnel.svg
[release-url]: https://github.com/snsinfu/reverse-tunnel/releases
[travis-badge]: https://travis-ci.org/snsinfu/reverse-tunnel.svg?branch=master
[travis-url]: https://travis-ci.org/snsinfu/reverse-tunnel
[license-badge]: https://img.shields.io/badge/license-MIT-blue.svg
[license-url]: https://raw.githubusercontent.com/snsinfu/reverse-tunnel/master/LICENSE
[report-badge]: https://goreportcard.com/badge/github.com/snsinfu/reverse-tunnel
[report-url]: https://goreportcard.com/report/github.com/snsifnu/reverse-tunnel

**rtun** is a tool for exposing TCP and UDP ports to the Internet via a public
gateway server. You can expose ssh and mosh server on a machine behind firewall
and NAT.

- [Build](#build)
- [Usage](#usage)
  - [Gateway server](#gateway-server)
  - [Agent](#agent)
- [License](#license)

## Build

Outside your GOPATH, run

```console
export GO111MODULE=on
git clone https://github.com/snsinfu/reverse-tunnel
cd reverse-tunnel
make
```

The `make` command produces two executable files: `rtun` and `rtun-server`. Put
`rtun-server` in a public server and `rtun` in a local machine.

## Usage

### Gateway server

Create a configuration file named `rtun-server.yml`:

```yaml
# Gateway server binds to this address to communicate with agents.
control_address: 0.0.0.0:9000

# List of authorized agents follows.
agents:
  - auth_key: a79a4c3ae4ecd33b7c078631d3424137ff332d7897ecd6e9ddee28df138a0064
    ports: [10022/tcp, 60000/udp]
```

You may want to generate `auth_key` with `openssl rand -hex 32`. Agents are
identified by their keys and the agents may only use the whitelisted ports
listed in the configuration file.

Then, start gateway server:

```console
./rtun-server
```

Now agents can connect to the gateway server and start reverse tunneling. The
server and agent uses WebSocket for communication, so the gateway server may be
placed behind an HTTPS reverse proxy like caddy. This way the tunnel can be
secured by TLS.

#### Standalone TLS

`rtun-server` supports automatic acquisition and renewal of TLS certificate.
Set control address to `:443` and `domain` to the domain
name of the public gateway server.

```
control_address: :443

lets_encrypt:
  domain: rtun.example.com
```

Non-root user can not use port 443 by default. You may probably want to allow
`rtun-server` bind to privileged port using `setcap` on Linux:

```
sudo setcap cap_net_bind_service=+ep rtun-server
```

### Agent

Create a configuration file named `rtun.yml`:

```yaml
# Specify the gateway server.
gateway_url: ws://the-gateway-server.example.com:9000

# A key registered in the gateway server configuration file.
auth_key: a79a4c3ae4ecd33b7c078631d3424137ff332d7897ecd6e9ddee28df138a0064

forwards:
  # Forward 10022/tcp on the gateway server to localhost:22 (tcp)
  - port: 10022/tcp
    destination: 127.0.0.1:22

  # Forward 60000/udp on the gateway server to localhost:60000 (udp)
  - port: 60000/udp
    destination: 127.0.0.1:60000
```

And run agent:

```console
./rtun
```

Note: When you are using TLS on the server the gateway URL should start with
`wss://` instead of `ws://`. In this case, the port number should most likely
be the default:

```yaml
gateway_url: wss://the-gateway-server.example.com
```

## License

MIT License.
