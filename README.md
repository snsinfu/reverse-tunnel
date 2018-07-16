Reverse tunnel TCP and UDP
==========================

[![Release][release-badge]][release-url]
[![Build Status][travis-badge]][travis-url]
[![MIT License][license-badge]][license-url]

[release-badge]: https://img.shields.io/github/release/snsinfu/reverse-tunnel.svg
[release-url]: https://raw.githubusercontent.com/snsinfu/reverse-tunnel/releases
[travis-badge]: https://travis-ci.org/snsinfu/reverse-tunnel.svg?branch=master
[travis-url]: https://travis-ci.org/snsinfu/reverse-tunnel
[license-badge]: https://img.shields.io/badge/license-MIT-blue.svg
[license-url]: https://raw.githubusercontent.com/snsinfu/reverse-tunnel/master/LICENSE
[report-badge]: https://goreportcard.com/badge/github.com/snsinfu/reverse-tunnel
[report-url]: https://goreportcard.com/report/github.com/snsifnu/reverse-tunnel

This repository contains **rtun**, a tool for easily exposing TCP and UDP ports
to the Internet via a public gateway server. It can be used, for example, to
expose ssh and mosh server behind firewall and NAT.

- [Build](#build)
- [Usage](#usage)
  - [Gateway server](#gateway-server)
  - [Agent](#agent)
- [License](#license)

## Build

```console
git clone https://github.com/snsinfu/reverse-tunnel
cd reverse-tunnel
make
```

This produces two executable files named `rtun` and `rtun-server`. Place
`rtun-server` in a public server.

## Usage

### Gateway server

Create a configuration file named `rtun-server.yml`:

```yaml
# Gateway server binds to this address to communicate with agents.
control_address: 0.0.0.0:9000

# List of authorized agents follows.
agents:
  - auth_key: a79a4c3ae4ecd33b7c078631d3424137ff332d7897ecd6e9ddee28df138a0064
    ports: [10022/tcp, 10022/udp]
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
    destination: localhost:22

  # Forward 10022/udp on the gateway server to localhost:10022 (udp)
  - port: 10022/udp
    destination: localhost:10022
```

And run agent:

```console
./rtun
```

Note: When you are using HTTPS reverse proxy the gateway URL should start with
`wss://` instead of `ws://`.

## License

MIT License.
