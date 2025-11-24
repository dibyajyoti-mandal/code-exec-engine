package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/dibyajyoti-mandal/code-exec-engine/constants"
	executor "github.com/dibyajyoti-mandal/code-exec-engine/exec"
)

type Job struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	Image    string `json:"image"`
}

var limiter = make(chan struct{}, constants.IMAGE_LIMIT) //buffered channel as a semaphore
var activeCount = 0
var mu sync.Mutex

// Fake job queue
var jobQueue = make(chan Job, constants.JQCHANNEL)

func enqueueTestJobs() {
	jobs := []Job{
		{Language: "python", Code: `print("Hello 1")`},
		{Language: "python", Code: `print("Hello 2")`},
		{Language: "cpp", Code: `
#include <iostream>
using namespace std;
int main(){
int n = 4;
for(int i=1; i<=n; i++){
	cout<<i<<" ";
}cout<<endl;

}`},
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

			mu.Lock()
			activeCount++
			fmt.Println("Active containers =", activeCount)
			mu.Unlock()

			go func(job Job) {
				fmt.Printf("[Worker %d] Starting job (%s)\n", workerID, job.Language)
				defer func() {
					<-limiter

					mu.Lock()
					activeCount--
					fmt.Println("Active containers =", activeCount)
					mu.Unlock()
				}()

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
