# Docker images

- [Server](#server)
- [Agent](#agent)


## Server

| Environment variable | Description                                  |
|----------------------|----------------------------------------------|
| `RTUN_AGENT`         | List of agents (allowed ports and auth keys) |
| `RTUN_TLS`           | (optional) Let's Encrypt domain              |

```sh
docker run -it \
  -p 9000:9000 \
  -e RTUN_AGENT="8080/tcp,30000/udp @ SampleKey-bfeeb1356a458eabef49e7e7" \
  snsinfu/rtun-server
```


## Agent

| Environment variable | Description                             |
|----------------------|-----------------------------------------|
| `RTUN_GATEWAY`       | WebSocket URL of rtun-server            |
| `RTUN_KEY`           | Auth key associated to the agent        |
| `RTUN_FORWARD`       | List of port forwardings                |

```sh
docker run -it \
  -e RTUN_GATEWAY="ws://0.1.2.3:9000" \
  -e RTUN_KEY="SampleKey-bfeeb1356a458eabef49e7e7" \
  -e RTUN_FORWARD="8080/tcp:192.168.1.10:8080" \
  snsinfu/rtun
```

docker-compose:

```yaml
version: "3.8"

services:
  web:
    image: nginx

  rtun:
    image: snsinfu/rtun
    environment:
      RTUN_GATEWAY: ws://0.1.2.3:9000
      RTUN_KEY: SampleKey-bfeeb1356a458eabef49e7e7
      RTUN_FORWARD: >
        8080/tcp:web:80
```
