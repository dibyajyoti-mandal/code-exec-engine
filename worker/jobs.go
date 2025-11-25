package main

import "github.com/dibyajyoti-mandal/code-exec-engine/models"

func EnqueueTestJobs() {
	jobs := []models.Job{
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
		{
			Language: "python",
			Code: `
while True:
    pass
`,
		},
		{
			Language: "python",
			Code: `
import time
time.sleep(2.5)
print("Done sleeping")
`,
		},
	}

	for _, j := range jobs {
		jobQueue <- j
	}
}
