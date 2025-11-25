package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/dibyajyoti-mandal/code-exec-engine/constants"
	executor "github.com/dibyajyoti-mandal/code-exec-engine/exec"
	"github.com/dibyajyoti-mandal/code-exec-engine/models"
)

var jobQueue = make(chan models.Job, constants.JQCHANNEL)

func workerLoop(workerID int) {
	fmt.Printf("[Worker %d] started\n", workerID)

	// Each worker: MAX 2 containers
	workerLimiter := make(chan struct{}, 2)
	var mu sync.Mutex
	active := 0

	for {
		select {
		case job := <-jobQueue:

			workerLimiter <- struct{}{}

			mu.Lock()
			active++
			fmt.Printf("[Worker %d] Active containers = %d\n", workerID, active)
			mu.Unlock()

			// Run job in worker goroutine
			go func(job models.Job) {
				defer func() {
					<-workerLimiter // release slot

					mu.Lock()
					active--
					fmt.Printf("[Worker %d] Active containers = %d\n", workerID, active)
					mu.Unlock()
				}()

				fmt.Printf("[Worker %d] Starting job (%s)\n", workerID, job.Language)
				runJob(workerID, job)

			}(job)

		default:
			time.Sleep(150 * time.Millisecond)
		}
	}
}

func runJob(workerID int, job models.Job) {
	var image string

	switch job.Language {
	case "python":
		image = "code-exec-python"
	case "cpp":
		image = "code-exec-cpp"
	default:
		fmt.Printf("[Worker %d] Unknown language: %s\n", workerID, job.Language)
		return
	}

	result := executor.RunInDocker(image, job.Code)
	out, _ := json.MarshalIndent(result, "", "  ")

	fmt.Printf("[Worker %d] Job finished:\n%s\n\n", workerID, string(out))
}

func main() {
	EnqueueTestJobs()

	// Spawn workers
	for i := 1; i <= 3; i++ {
		go workerLoop(i)
	}

	select {} // keep running
}

/*to do - graceful shutdown:
To achieve this, we need two things-

sync.WaitGroup: To track active jobs and ensure we don't exit until the counter reaches zero.

Signal Handling: To intercept Ctrl+C so we can close the queue safely instead of crashing immediately.
*/
