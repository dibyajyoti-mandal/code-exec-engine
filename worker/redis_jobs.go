package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/dibyajyoti-mandal/code-exec-engine/constants"
	"github.com/dibyajyoti-mandal/code-exec-engine/models"
	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Println(">> Connected to Redis")

	// MKSTREAM ensures the stream is created if it doesn't exist
	// "$" means "start reading from now" (don't replay old messages on startup)
	err := rdb.XGroupCreateMkStream(ctx, constants.STREAM_KEY, constants.GROUP_NAME, "$").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		log.Printf("Error creating consumer group: %v", err)
	}
}

func StartRedisConsumer() {
	fmt.Println(">> Redis Consumer Started. Listening for 'code-jobs'...")

	for {
		//Read from Stream
		entries, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    constants.GROUP_NAME,
			Consumer: constants.CONSUMER,
			Streams:  []string{constants.STREAM_KEY, ">"},
			Count:    1,
			Block:    0,
		}).Result()

		if err != nil {
			log.Printf("Redis Read Error: %v", err)
			continue
		}

		for _, stream := range entries {
			for _, message := range stream.Messages {

				// Extract the JSON string
				val, ok := message.Values["job"]
				if !ok {
					log.Println("Error: Key 'job' not found in Redis message")
					rdb.XAck(ctx, constants.STREAM_KEY, constants.GROUP_NAME, message.ID)
					continue
				}

				jsonStr, ok := val.(string)
				if !ok {
					log.Println("Error: 'job' value is not a string")
					continue
				}

				// Unmarshal into Job Struct
				var job models.Job
				err := json.Unmarshal([]byte(jsonStr), &job)
				if err != nil {
					log.Printf("Error unmarshaling job: %v", err)
					// Still Ack to prevent infinite retry loop on bad data
					rdb.XAck(ctx, constants.STREAM_KEY, constants.GROUP_NAME, message.ID)
					continue
				}

				// Send to Internal Queue
				fmt.Printf(">> [Redis] Received Job: %s (Client: %s)\n", job.ID, job.ClientID)
				jobQueue <- job

				// Acknowledge
				rdb.XAck(ctx, constants.STREAM_KEY, constants.GROUP_NAME, message.ID)
			}
		}
	}
}
