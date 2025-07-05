# Project Brief: TON Indexer Health-Checker

## 1. High-Level Objective

The primary goal is to develop a production-ready health-checker utility in Go. This tool is designed to monitor a critical component of the TON blockchain infrastructure, the `index-worker`. It will be integrated into a Docker Swarm orchestration system to enable automatic recovery of the service.

## 2. Core Problem

The `index-worker` service, which indexes TON blockchain data, has a critical flaw: it can enter a "hung" state where the process remains active but stops processing new blocks. Standard Docker health checks (`restart: always`) are ineffective because the process does not crash. This leads to data on the marketplace becoming stale and requires manual intervention.

## 3. Proposed Solution

The solution is to implement an external, intelligent health check that verifies the *actual* operational status of the indexer, not just its process state. This will be achieved by a Go utility that:

1.  Connects to a Redis instance where the `index-worker` leaves a timestamp of its last successful operation.
2.  Compares this timestamp against the current time.
3.  Exits with a specific status code (`0` for healthy, `1` for unhealthy) based on a configurable time threshold.

Docker Swarm will use this utility's exit code to determine the container's health and automatically restart it after a series of consecutive failures, thus automating the recovery process.
