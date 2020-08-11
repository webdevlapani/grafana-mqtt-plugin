# Grafana MQTT Datasource Plugin for Tvarit Cloud Platform

[![CircleCI](https://circleci.com/gh/tvarit-foggy/grafana-mqtt-plugin/tree/master.svg?style=svg)](https://circleci.com/gh/tvarit-foggy/grafana-mqtt-plugin/tree/master)

![Architecture Diagram](https://github.com/tvarit-foggy/grafana-mqtt-plugin/raw/master/src/img/architecture.png)

Landing zone ingests IoT data from MQTT and stores it in S3 bucket under the key <topic>/<timestamp_ms>. Topic name and client ID must start with <orgId>/<datasourceId>/ which are provided during configuration.

### Installation
To install, download zip file from release page (for stable version) or download repository as zip file (for installing from git master)
```
grafana-cli --pluginUrl <path_to_zip_file> plugins install tvarit-mqtt-datasource
```

The plugin is not signed. To enable this plugin, change the following

```
[plugins]
...
# Enter a comma-separated list of plugin identifiers to identify plugins that are allowed to be loaded even if they lack a valid signature.
allow_loading_unsigned_plugins =
```
to

```
[plugins]
...
# Enter a comma-separated list of plugin identifiers to identify plugins that are allowed to be loaded even if they lack a valid signature.
allow_loading_unsigned_plugins = tvarit-mqtt-datasource
```

## Getting started
1. Login to Tvarit Cloud Platform and navigate to Configuration > Datasources.
2. Click on add data source, select MQTT Datasource under Tvarit GmbH and give it a name.
3. The page now shows MQTT host and port details. These should be configured into MQTT client.
4. Client certificates must be registered with Tvarit before a client can communicate with MQTT Datasource. For this tutorial, we will use one-click certificate creation. This will generate a certificate, public key, and private key. Click on register certificate and click on one-click certificate creation.
5. Download the certificate files and save them in a safe place. These files cannot be retrieved after you close this page. You will also see topic and client id prefix mentioned on this page. Messages can be published to topics and client ids wit this prefix using this certificate.

## Development

### Common

1. Fork this repository.
2. Create a directory called `grafana-plugins` in your preferred workspace.
3. Clone the fork in `grafana-plugins` directory.
4. Find the `plugins` property in the Grafana configuration file and set the `plugins` property to the path of your `grafana-plugins` directory. Refer to the [Grafana configuration documentation](https://grafana.com/docs/grafana/latest/installation/configuration/#plugins) for more information.
```INI
[paths]
plugins = "/path/to/grafana-plugins"
```
5. Fill relevant values in `config.ini`.
```INI
service = Cloud
zone = Landing Zone
environment = alpha

[aws]
account_id =
access_key =
secret_key =
```
6. Restart Grafana if it’s already running.

### Frontend

1. Install prerequisites
    * NodeJS >=12,<13 (https://nodejs.org/en/download/)
    * yarn (https://classic.yarnpkg.com/en/docs/install/#debian-stable)

2. Install dependencies
```BASH
yarn install
```

3. Build plugin in development mode or run in watch mode
```BASH
yarn dev
```
or
```BASH
yarn watch
```

4. Build plugin in production mode
```BASH
yarn build
```

### Backend

1. Install prerequisites
    * Go 1.14.* (https://golang.org/dl/)
    * Mage (https://magefile.org/)

2. Update [Grafana plugin SDK for Go](https://grafana.com/docs/grafana/latest/developers/plugins/backend/grafana-plugin-sdk-for-go/) dependency to the latest minor version:

```bash
go get -u github.com/grafana/grafana-plugin-sdk-go
```

3. Build backend plugin binary for Linux in development mode:
```BASH
./scripts/go/bin/bra run
```
or
```BASH
mage -v Trace
```

4. Build backend plugin binaries for Linux, Windows and Darwin:
```BASH
mage -v
```

5. List all available Mage targets for additional commands:
```BASH
mage -l
```
