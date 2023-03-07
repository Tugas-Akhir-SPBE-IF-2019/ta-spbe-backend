# ta-spbe-backend

## Run Locally

Note: This project use [docker](https://github.com/docker/cli) and [docker-compose](https://github.com/docker/compose) to run all required infrastructure objects locally.

After git clone, use following commands sequence:

- `make vendor`. To download Go module dependencies into `vendor` folder. (first time only, or on each code dependency changes).
- `make config`. To copy `config.toml.example` to `config.toml` (this file is used to set several settings (database, monitoring tools, etc.) for the server to run).
  This is for local reference, and will be `exported` right before running the server (see [docker-compose.yml](./docker-compose.yml)). (first time only, or on each config changes).
- `make server/start`. To start the servers (infrastructure and API server).