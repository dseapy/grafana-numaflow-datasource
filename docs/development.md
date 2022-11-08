## Requirements
* Grafana 9.2+
* Go 1.19+
* NodeJS 16
* Yarn

For development, set the following in Grafana's `grafana.ini`:
```ini
# Allows for not needing to sign the plugin
app_mode = development

# Can be different directory, but make sure the directory exists
[paths]
plugins = /var/lib/grafana/plugins

# Will show logs for the backend datasource code in the grafana logs
[log]
level = debug
```

For more details on building grafana data-sources see [build-a-data-source-plugin](https://grafana.com/tutorials/build-a-data-source-plugin/).

## Frontend

Install dependencies
```bash
yarn install
```

Build
```bash
yarn build
```

## Backend

Build:
```bash
# creates binaries for Linux, Windows and Darwin in "dist" directory
mage -v
```

## Install
Ensure `dist/*` are available at `<grafana-plugins-dir>/numaflow/*`, with `<grafana-plugins-dir>` being
what is configured in the `grafana.ini` under `paths.plugins` (see `Requirements` above).

Restart grafana and datasource should be updated.
