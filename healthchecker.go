package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// IndexerState represents the structure of the JSON object stored in Redis.
type IndexerState struct {
	GenUtime int64 `json:"gen_utime"`
}

// getConfig reads a configuration value from an environment variable,
// falling back to a default value if the variable is not set.
func getConfig(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	// --- Configuration ---
	redisHost := getConfig("REDIS_HOST", "event-cache")
	redisPort := getConfig("REDIS_PORT", "6379")
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	stateKey := getConfig("REDIS_STATE_KEY", "last_mc_seqno")
	maxDelayStr := getConfig("MAX_DELAY_SECONDS", "300")

	maxDelay, err := strconv.ParseInt(maxDelayStr, 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FAIL: Invalid MAX_DELAY_SECONDS value: %v\n", err)
		os.Exit(1)
	}

	// --- Health Check Logic ---
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	defer rdb.Close()

	stateJSON, err := rdb.Get(ctx, stateKey).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Fprintf(os.Stderr, "FAIL: State key '%s' not found in Redis.\n", stateKey)
		} else {
			fmt.Fprintf(os.Stderr, "FAIL: Cannot get state from Redis: %v\n", err)
		}
		os.Exit(1)
	}

	var state IndexerState
	if err := json.Unmarshal([]byte(stateJSON), &state); err != nil {
		fmt.Fprintf(os.Stderr, "FAIL: Cannot parse JSON state from key '%s': %v\n", stateKey, err)
		os.Exit(1)
	}

	if state.GenUtime == 0 {
		fmt.Fprintf(os.Stderr, "FAIL: 'gen_utime' is missing or zero in state from key '%s'.\n", stateKey)
		os.Exit(1)
	}

	delay := time.Now().Unix() - state.GenUtime
	if delay < 0 {
		// This can happen if clocks are out of sync. Treat as a failure.
		fmt.Fprintf(os.Stderr, "FAIL: System clock seems to be behind the indexer's clock (delay: %d). Check time synchronization.\n", delay)
		os.Exit(1)
	}
	
	if delay > maxDelay {
		fmt.Fprintf(os.Stderr, "FAIL: Indexer delay is %d seconds (limit %d).\n", delay, maxDelay)
		os.Exit(1)
	}

	fmt.Printf("OK: Indexer delay is %d seconds.\n", delay)
	os.Exit(0)
}
