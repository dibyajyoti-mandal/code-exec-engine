package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/dibyajyoti-mandal/code-exec-engine/constants"
	executor "github.com/dibyajyoti-mandal/code-exec-engine/exec"
	"github.com/dibyajyoti-mandal/code-exec-engine/models"
	"github.com/dibyajyoti-mandal/code-exec-engine/socket"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	jobQueue    = make(chan models.Job, constants.JQCHANNEL)
	resultQueue = make(chan models.Result, constants.JQCHANNEL)
	upgrader    = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func main() {

	//initalize redis conn
	InitRedis()

	//fire workers

	for i := 1; i <= constants.WORKER_COUNT; i++ {
		go workerLoop(i)
	}
	go resultBroadcaster()

	go StartRedisConsumer()

	http.HandleFunc("/ws", handleWebSocket)

	fmt.Println("Code Execution Engine Running on", constants.SERVER_PORT)
	if err := http.ListenAndServe(constants.SERVER_PORT, nil); err != nil {
		panic(err)
	}
}

func DownloadTestcase(client *s3.Client, bucket, key, localPath string) error {
	out, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return err
	}
	defer out.Body.Close()

	f, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, out.Body)
	return err
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	clientID := uuid.New().String()
	socket.Pool.Add(clientID, conn)
	fmt.Printf("[WS] Client connected: %s\n", clientID)

	defer socket.Pool.Remove(clientID)

	for {
		var job models.Job
		err := conn.ReadJSON(&job)
		if err != nil {
			break
		}

		job.ClientID = clientID
		if job.ID == "" {
			job.ID = uuid.New().String()
		}

		fmt.Printf("[WS] Job %s queued for Client %s\n", job.ID, clientID)
		jobQueue <- job
	}
}

func workerLoop(workerID int) {
	workerLimiter := make(chan struct{}, 2)
	var mu sync.Mutex
	active := 0

	for job := range jobQueue {
		workerLimiter <- struct{}{}

		mu.Lock()
		active++
		mu.Unlock()

		go func(j models.Job) {
			defer func() {
				<-workerLimiter
				mu.Lock()
				active--
				mu.Unlock()
			}()

			//processing update
			resultQueue <- models.Result{
				JobID:    j.ID,
				ClientID: j.ClientID,
				Status:   "Processing",
			}

			//run job
			res := runJobLogic(workerID, j)
			res.Status = "Completed"
			resultQueue <- res
		}(job)
	}
}

func runJobLogic(workerID int, job models.Job) models.Result {
	var image string
	switch job.Language {
	case "python":
		image = "code-exec-python"
	case "cpp":
		image = "code-exec-cpp"
	default:
		return models.Result{
			JobID: job.ID, ClientID: job.ClientID,
			Error: "Unsupported Language",
		}
	}
	_ = executor.MarkProblemDone()
	executionResult := executor.RunInDocker(image, job.Code)

	executionResult.JobID = job.ID
	executionResult.ClientID = job.ClientID

	return executionResult
}

func resultBroadcaster() {
	for res := range resultQueue {

		if res.Status == "Completed" {
			out, _ := json.MarshalIndent(res, "", "  ")
			fmt.Printf("\n[Terminal] Job %s Finished:\n%s\n\n", res.JobID, string(out))
		} else {
			fmt.Printf("[Terminal] Job %s is Processing...\n", res.JobID)
		}

		if res.ClientID == "BROADCAST" {
			socket.Pool.Broadcast(res)
			fmt.Println("[Global] Broadcast on", res.JobID, res.Status)
		} else {
			socket.Pool.SendResult(res.ClientID, res)
			if res.Status == "Completed" {
				fmt.Printf("[Hub] Result routed to Client %s\n", res.ClientID)
			}
		}
	}
}

// run: "docker build -t <Image name> ." for each language
