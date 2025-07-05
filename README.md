# TON Indexer Health Checker

A lightweight, production-ready Go utility designed to monitor the health of a TON `index-worker` service within a Docker Swarm environment.

## 1. Purpose

The `index-worker` component in the TON blockchain indexing pipeline can sometimes enter a "hung" state where it stops processing new blocks but the process itself does not crash. Standard Docker restart policies are ineffective against this type of silent failure.

This health checker provides an intelligent, external check to determine if the indexer is *actually* working. It does this by checking a timestamp that the indexer writes to a Redis cache. If the timestamp is too old, the health checker exits with a failure code, signaling to Docker Swarm that the container is unhealthy and needs to be restarted.

This enables fully automated recovery for a critical piece of infrastructure.

## 2. How It Works

1.  **Docker Swarm** executes this `healthchecker` utility inside the `index-worker` container on a schedule.
2.  The utility connects to a configured **Redis** instance.
3.  It reads a specific key (`last_mc_seqno` by default) which contains a JSON payload with a `gen_utime` (Unix timestamp) field.
4.  It calculates the delay between the current time and the `gen_utime`.
5.  **If the delay is within the allowed limit:** It prints an `OK` message to `stdout` and exits with code `0`.
6.  **If the delay exceeds the limit:** It prints a `FAIL` message to `stderr` and exits with code `1`.

Docker Swarm uses these exit codes to manage the container's lifecycle, automatically restarting it after a configured number of consecutive failures.

## 3. Configuration

The utility is configured entirely through environment variables, making it ideal for containerized deployments.

| Variable            | Description                                       | Default         |
| ------------------- | ------------------------------------------------- | --------------- |
| `REDIS_HOST`        | The hostname of the Redis server.                 | `event-cache`   |
| `REDIS_PORT`        | The port of the Redis server.                     | `6379`          |
| `REDIS_STATE_KEY`   | The key in Redis holding the indexer state JSON.  | `last_mc_seqno` |
| `MAX_DELAY_SECONDS` | Max allowed delay in seconds before failure.      | `300`           |

## 4. Building the Utility

To build the binary from the source code, run the following command in the project root:

```sh
go build -o healthchecker .
```

This will create a `healthchecker` executable in the current directory.

## 5. Docker Integration

This utility is designed to be used as a `healthcheck` in a Docker environment.

### Multi-stage Dockerfile

To add the health checker to your existing `index-worker` image without bloating it, use a multi-stage build.

**Modify your `Dockerfile`:**

```dockerfile
# === Stage 1: Build the healthchecker ===
FROM golang:1.21-alpine as builder

WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code and build
COPY . .
RUN CGO_ENABLED=0 go build -o /healthchecker .


# === Stage 2: Final application image ===
# Replace this with your actual base image
FROM toncenter/ton-indexer-base:latest

# ... (all your existing Dockerfile instructions for the index-worker) ...

# Copy the compiled healthchecker binary from the builder stage
COPY --from=builder /healthchecker /usr/local/bin/healthchecker

# Ensure the binary is executable
RUN chmod +x /usr/local/bin/healthchecker
```

### Docker Compose Healthcheck

To configure Docker Swarm to use the health checker, add the `healthcheck` section to your `index-worker` service in your `docker-compose.yaml` file.

**Example `docker-compose.yaml` service definition:**

```yaml
services:
  index-worker:
    image: your-custom-indexer-image:latest
    # ... other service configuration ...
    environment:
      - REDIS_HOST=my-redis-instance
      - MAX_DELAY_SECONDS=600
    healthcheck:
      # Command to execute for the health check
      test: ["CMD", "/usr/local/bin/healthchecker"]
      
      # How often to run the check
      interval: 60s
      
      # How long to wait for the check to complete
      timeout: 15s
      
      # Number of consecutive failures to be considered "unhealthy"
      retries: 3
      
      # Grace period for the container to start up before checks begin
      start_period: 3m
