# grafana_loki
A simple base config files to start up log gathering system, running on docker-compose, with custom config and persist data.

## step
- clone to remote machine.
    ```
    git clone https://github.com/hellodudu/grafana_loki.git
    ```

- make sure docker and docker-compose installed.

- check docker-log-driver [loki-docker-driver](https://github.com/grafana/loki/tree/master/cmd/docker-driver) installed.

- run with docker-compose.
    ```
    docker-compose up -d
    ```

- set datasource `http://loki:3100` in grafana.

- push logs with api `http://{public_ip}:3100/api/prom/push`.

## config and persist
- you can custom grafana and loki config in `./config/grafana.ini` and `./config/loki/loki-local-config.yaml`.
- persist data saved in `./data`.
- both config and persist files are mapped to docker container, you can config them in `./docker-compose.yaml`
