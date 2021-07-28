# Changelog

## Unreleased

- Added basic config checks.
- Fixed agent to gracefully close tunnels on interrupt. This should fix
  dangling connections on server upon stopping agent.

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
