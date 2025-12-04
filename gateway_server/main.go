package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dibyajyoti-mandal/code-exec-engine/gateway/constants"
	"github.com/dibyajyoti-mandal/code-exec-engine/gateway/models"
	"github.com/dibyajyoti-mandal/code-exec-engine/gateway/queue"
	"github.com/google/uuid"
)

var publisher queue.Publisher

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", constants.FRONTEND)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	// choose queue backend
	publisher = queue.NewRedisPublisher(constants.REDIS_ADDR)

	http.HandleFunc("/health/redis", enableCORS(handleRedisHealth))

	http.HandleFunc("/submit", enableCORS(handleSubmit))

	log.Println("Gateway server running on", constants.SERVER_PORT)
	log.Fatal(http.ListenAndServe(constants.SERVER_PORT, nil))
}

func handleRedisHealth(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// try basic PING
	redisPub, ok := publisher.(*queue.RedisPublisher)
	if !ok {
		http.Error(w, "Redis publisher not enabled", http.StatusBadRequest)
		return
	}

	_, err := redisPub.Client().Ping(ctx).Result()
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"redis":  "connected",
	})
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		ClientID string `json:"client_id"`
		Language string `json:"language"`
		Code     string `json:"code"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	job := models.Job{
		ID:        uuid.NewString(),
		ClientID:  payload.ClientID,
		Language:  payload.Language,
		Code:      payload.Code,
		Timestamp: time.Now().Unix(),
	}

	err = publisher.Publish(job)
	if err != nil {
		http.Error(w, "failed to enqueue job", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "queued",
		"jobId":  job.ID,
	})
}
