# SystemD Service UI

`systemd-service-ui` provides a web user interface to control a subset of systemd services.

## Building

`systemd-service-ui` must be built as a binary. While it can be built as a container, it needs access to the host dbus to control system services and containers.

To build the RPM, run:

```sh
./scripts/release.sh systemd-service-ui
```

This will compile the binary and build an RPM package in build/dist.

## Running

The RPM package includes a Systemd service called `systemd-service-ui`, this means the service can be controlled with:

```sh
systemctl start systemd-service-ui
systemctl stop systemd-service-ui
systemctl restart systemd-service-ui
```

The default path for the configuration file is `/etc/systemd-service-ui.yml`.
