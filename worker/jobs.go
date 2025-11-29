package main

import (
	"fmt"
	"time"

	"github.com/dibyajyoti-mandal/code-exec-engine/models"
	"github.com/google/uuid"
)

func EnqueueTestJobs() {
	fmt.Println(">> System Test Jobs will run in 5 seconds (Connect WS now!)...")
	time.Sleep(5 * time.Second)

	jobs := []models.Job{
		{
			ID:       uuid.New().String(),
			ClientID: "BROADCAST",
			Language: "python",
			Code:     `print("Hello 1")`,
		},
		{
			ID:       uuid.New().String(),
			ClientID: "BROADCAST",
			Language: "python",
			Code:     `print("Hello 2")`,
		},
		{
			ID:       uuid.New().String(),
			ClientID: "BROADCAST",
			Language: "cpp",
			Code: `
#include <iostream>
using namespace std;
int main(){
    int n = 4;
    for(int i=1; i<=n; i++){
        cout<<i<<" ";
    }
    cout<<endl;
    return 0;
}`,
		},
		{
			ID:       uuid.New().String(),
			ClientID: "BROADCAST",
			Language: "python",
			Code: `
while True:
    pass
`,
		},
		{
			ID:       uuid.New().String(),
			ClientID: "BROADCAST",
			Language: "python",
			Code: `
import time
time.sleep(4.5)
print("Done sleeping")
`,
		},
	}

	for _, j := range jobs {
		jobQueue <- j
		fmt.Printf(">> Enqueued System Job: %s\n", j.ID)
	}
}
