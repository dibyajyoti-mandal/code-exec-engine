package main

import (
	"encoding/json"
	"fmt"
	"time"

	executor "github.com/dibyajyoti-mandal/code-exec-engine/exec"
)

type Job struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	Image    string `json:"image"`
}

var limiter = make(chan struct{}, 2)

// Fake job queue
var jobQueue = make(chan Job, 10)

func enqueueTestJobs() {
	jobs := []Job{
		{Language: "python", Code: `print("Hello 1")`},
		{Language: "python", Code: `print("Hello 2")`},
		{Language: "cpp", Code: `#include <iostream>
int main(){ std::cout << "Hello from C++"; }`},
	}

	for _, j := range jobs {
		jobQueue <- j
	}
}

func workerLoop(workerID int) {
	fmt.Println("[Worker", workerID, "] started")

	for {
		select {
		case job := <-jobQueue:
			limiter <- struct{}{}

			go func(job Job) {
				fmt.Printf("[Worker %d] Starting job (%s)\n", workerID, job.Language)
				defer func() { <-limiter }()

				var image string
				switch job.Language {
				case "python":
					image = "code-exec-python"
				case "cpp":
					image = "code-exec-cpp"
				}

				result := executor.RunInDocker(image, job.Code)
				out, _ := json.MarshalIndent(result, "", "  ")
				fmt.Printf("[Worker %d] Job finished:\n%s\n", workerID, string(out))
			}(job)

		default:
			// Idle wait → prevents deadlock
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func main() {
	enqueueTestJobs()

	for i := 1; i <= 3; i++ {
		go workerLoop(i)
	}

	select {} // keep program alive forever
}
