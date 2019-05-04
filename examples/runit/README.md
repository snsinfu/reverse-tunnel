# [Runit](http://smarden.org/runit/) service files

## Agent

1. Make `~/service` directory
2. Copy `rtun` directory into `~/service/`
3. Copy `rtun` binary into `~/service/rtun/`
4. Configure `~/service/rtun/rtun.yml`
5. Run `runsvdir ~/service`
6. Done

## Server

1. Make `~/service` directory
2. Copy `rtun-server` directory into `~/service/`
3. Copy `rtun-server` binary into `~/service/rtun-server/`
4. Configure `~/service/rtun-server/rtun-server.yml`
5. Run `runsvdir ~/service`
6. Done

## Start runsvdir at boot time

You may want to automate the execution of `runsvdir` by registering the command
to crontab. The crontab line would look like this:

```
@reboot daemonize runsvdir ~/service
```
