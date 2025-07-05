# Technical Context: TON Indexer Health-Checker

## 1. Core Technology Stack

- **Language:** Go (version 1.19 or higher)
- **Primary Dependency:** `github.com/go-redis/redis/v8` for Redis communication.
- **Build Environment:** Docker (using a multi-stage build).
- **Orchestration:** Docker Swarm.

## 2. Go Utility (`healthchecker`) Requirements

- **Binary Name:** The compiled executable must be named `healthchecker`.
- **Functionality:**
    - It is a single-execution command-line utility. It must not contain any long-running loops.
    - It must read its configuration from environment variables.
    - It must connect to Redis, fetch a specific key, parse the JSON content, and check the `gen_utime` timestamp.
    - It must exit with code `0` on success and `1` on failure.
    - Success messages must be printed to `stdout`.
    - Failure messages must be printed to `stderr`.

## 3. Configuration via Environment Variables

The utility must support the following environment variables, with sensible defaults:

- `REDIS_HOST`: The hostname of the Redis server.
  - **Default:** `event-cache`
- `REDIS_PORT`: The port of the Redis server.
  - **Default:** `6379`
- `REDIS_STATE_KEY`: The key in Redis that holds the indexer's state.
  - **Default:** `last_mc_seqno`
- `MAX_DELAY_SECONDS`: The maximum allowed delay in seconds before the indexer is considered unhealthy.
  - **Default:** `300`

## 4. Deployment & Integration

- **Dockerfile:** A multi-stage `Dockerfile` is required.
    - **Stage 1 (builder):** Uses a `golang` base image to build the `healthchecker` binary.
    - **Stage 2 (final):** Uses the existing `toncenter/ton-indexer-base:latest` image and copies the compiled binary from the `builder` stage into `/usr/local/bin/`.
- **Docker Compose:** The `docker-compose.yaml` for the `index-worker` service must be updated with a `healthcheck` section that specifies how to run the utility.
    - **Test Command:** `["CMD", "/usr/local/bin/healthchecker"]`
    - **Parameters:** `interval`, `timeout`, `retries`, and `start_period` must be configured appropriately.
