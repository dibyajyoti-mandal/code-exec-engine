package queue

import (
	"context"
	"encoding/json"
	"log"

	"github.com/dibyajyoti-mandal/code-exec-engine/gateway/constants"
	"github.com/dibyajyoti-mandal/code-exec-engine/gateway/models"
	"github.com/redis/go-redis/v9"
)

type RedisPublisher struct {
	client *redis.Client
	stream string
}

func NewRedisPublisher(addr string) *RedisPublisher {
	return &RedisPublisher{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
		stream: constants.REDIS_STREAM,
	}
}

func (r *RedisPublisher) Publish(job models.Job) error {
	body, _ := json.Marshal(job)

	_, err := r.client.XAdd(context.Background(), &redis.XAddArgs{
		Stream: r.stream,
		Values: map[string]interface{}{
			"job": string(body),
		},
	}).Result()

	if err != nil {
		log.Println("Redis publish error:", err)
	}

	return err
}

func (r *RedisPublisher) Client() *redis.Client {
	return r.client
}
