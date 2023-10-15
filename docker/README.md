# Docker images

- [Server](#server)
- [Agent](#agent)
  - [Docker-compose](#docker-compose)


## Server

| Environment variable | Description                                  |
|----------------------|----------------------------------------------|
| `RTUN_AGENT`         | List of agents (allowed ports and auth keys) |
| `RTUN_PORT`          | The port the server binds to (default: 9000) |
| `RTUN_TLS_CERT`      | TLS cert path (optional)                     |
| `RTUN_TLS_KEY`       | TLS key path (optional)                      |

```sh
docker run -it \
  -p 9000:9000 \
  -p 8080:8080 \
  -p 30000:30000 \
  -e RTUN_AGENT="8080/tcp,30000/udp @ samplebfeeb1356a458eabef49e7e7" \
  snsinfu/rtun-server
```

The `RTUN_AGENT` environment variable is a semicolon-separated list of rtun
agents to allow tunneling:

```sh
# Spaces and newlines are ignored.
RTUN_AGENT="
8080/tcp,30000/udp @ samplebfeeb1356a458eabef49e7e7;
9080/tcp @ sample1c96d79336ed361620d48d3e;
9090/tcp @ sampled915f77e1410fe92b62a435a"
```

Each agent specification should look like `ports @ key`. `ports` is a
comma-separated list of internet-facing ports to expose. `key` is a token
string used to authenticate the agent. In the above example, the first agent
uses key `samplebfeeb1356a458eabef49e7e7` to tunnel and expose TCP port 8080
and UDP port 30000 through the rtun-server.


## Agent

| Environment variable | Description                      |
|----------------------|----------------------------------|
| `RTUN_GATEWAY`       | WebSocket URL of rtun-server     |
| `RTUN_KEY`           | Auth key associated to the agent |
| `RTUN_FORWARD`       | List of port forwardings         |

```sh
docker run -it --network host \
  -e RTUN_GATEWAY="ws://0.1.2.3:9000" \
  -e RTUN_KEY="samplebfeeb1356a458eabef49e7e7" \
  -e RTUN_FORWARD="8080/tcp:localhost:8080" \
  snsinfu/rtun
```

The `--network host` option is required to forward to localhost. The
`RTUN_GATEWAY` environment variable specifies the WebSocket URL (ws:// or
wss://) of the `rtun-server` to use. The `RTUN_KEY` environment variable
specifies the authentication key to use. The `RTUN_FORWARD` environment
variable specifies tunnels as a comma-separated list:

```sh
RTUN_FORWARD="8080/tcp:localhost:8080, 30000/udp:192.168.1.10:30000"
```

Each tunnel rule looks like `internet-port/protocol:host:port`. In the above
example, the first rule specifies that the internet port `8080/tcp` on the
server should be tunneled to `localhost:8080`.


### Docker-compose

The following docker-compose example exposes local `nginx` as `http://0.1.2.3:8080`
by tunneling http connections through `rtun-server` running on a public server
0.1.2.3 on port 9000.

```yaml
version: "3.8"

services:
  web:
    image: nginx

  rtun:
    image: snsinfu/rtun
    environment:
      RTUN_GATEWAY: ws://0.1.2.3:9000
      RTUN_KEY: samplebfeeb1356a458eabef49e7e7
      RTUN_FORWARD: 8080/tcp:web:80
```
