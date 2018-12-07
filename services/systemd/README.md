# SystemD service unit files

## Agent

To make `rtun` auto-start at boot time:

1. Copy rtun binary to ~/bin
2. Put agent configurations into ~/.rtun.yml
3. Run these commands:
   ```
   $ mkdir -p ~/.config/systemd/user
   $ cp rtun.service ~/.config/systemd/user
   $ systemctl enable rtun --user
   $ loginctl enable-linger
   ```
4. Done

## Server

To make `rtun-server` auto-start at boot time:

1. Copy rtun-server binary to ~/bin
2. Put server configurations into ~/.rtun-server.yml
3. Run these commands:
   ```
   $ mkdir -p ~/.config/systemd/user
   $ cp rtun-server.service ~/.config/systemd/user
   $ systemctl enable rtun-server --user
   $ loginctl enable-linger
   ```
4. Done
