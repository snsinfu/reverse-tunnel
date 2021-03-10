Reverse tunnel TCP and UDP
==========================

[![Build Status][build-badge]][build-url]
[![Release][release-badge]][release-url]
[![MIT License][license-badge]](LICENSE.txt)

[build-badge]: https://github.com/snsinfu/reverse-tunnel/workflows/test/badge.svg
[build-url]: https://github.com/snsinfu/reverse-tunnel/actions?query=workflow%3Atest
[release-badge]: https://img.shields.io/github/release/snsinfu/reverse-tunnel.svg
[release-url]: https://github.com/snsinfu/reverse-tunnel/releases
[license-badge]: https://img.shields.io/badge/license-MIT-blue.svg

**rtun** is a tool for exposing TCP and UDP ports to the Internet via a public
gateway server. You can expose ssh and mosh server on a machine behind firewall
and NAT.

- [Build](#build)
- [Docker](#docker)
- [Usage](#usage)
  - [Gateway server](#gateway-server)
  - [Agent](#agent)
- [License](#license)


## Build

Compiled binaries are available in the [release page][release-url]. To build
your own ones, clone the repository and make:

```console
$ git clone https://github.com/snsinfu/reverse-tunnel
$ cd reverse-tunnel
$ make
```

Or,

```console
$ go build -o rtun github.com/snsinfu/reverse-tunnel/agent/cmd
$ go build -o rtun-server github.com/snsinfu/reverse-tunnel/server/cmd
```


## Docker

Docker images are available:

- https://hub.docker.com/r/snsinfu/rtun
- https://hub.docker.com/r/snsinfu/rtun-server

Quick usage:

```console
$ docker run -it \
  -p 8080:8080 -p 9000:9000 \
  -e RTUN_AGENT="8080/tcp @ samplebfeeb1356a458eabef49e7e7" \
  snsinfu/rtun-server

$ docker run -it --network host \
  -e RTUN_GATEWAY="ws://0.1.2.3:9000" \
  -e RTUN_KEY="samplebfeeb1356a458eabef49e7e7" \
  -e RTUN_FORWARD="8080/tcp:localhost:8080" \
  snsinfu/rtun
```

See [docker image readme](docker/README.md).


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
$ ./rtun-server
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
$ ./rtun
```

Note: When you are using TLS on the server the gateway URL should start with
`wss://` instead of `ws://`. In this case, the port number should most likely
be the default:

```yaml
gateway_url: wss://the-gateway-server.example.com
```


## License

MIT License.
