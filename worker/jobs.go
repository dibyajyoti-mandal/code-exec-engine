package main

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
		{Language: "python",
			Code: `
		while True:
		    pass
		`,
		},
		{
			Language: "python",
			Code: `
import time
time.sleep(3.5)
print("Done sleeping")
`,
		},
	}

	for _, j := range jobs {
		jobQueue <- j
	}
}
