# Changelog

## v1.4.0

- Replace Let's Encrypt by simple TLS certificates via `RTUN_TLS_CERT`, `RTUN_TLS_KEY`
- Added TLS dependencies in the container

## v1.3.2

- Added basic config checks.
- Fixed dangling connections on server upon stopping agent by interrupt.

## v1.3.1

- Fixed "websocket error: close 1006 (abnormal closure): unexpected EOF" on
  tunneling a lot of connections.
- Improved log messages.

## v1.3.0

- Fixed rtun-server to detect connection loss. This fix mitigates the annoying
  `"websocket: close 1000 (normal)" - recovering... ` loop.

## v1.2.3

- Fixed misconfiguration in the agent docker image

## v1.2.1

- Fixed EOF errors on TCP connection termination
- Bumped dependencies to the latest
- Added dockerfiles
