package executor

import (
	"bytes"
	"context"
	"os/exec"
	"time"
)

type Result struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Error  string `json:"error"`
}

func RunInDocker(image string, code string) Result {
	var stdout, stderr bytes.Buffer

	// timeout = 1 second
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		"docker", "run", "--rm",
		"--memory=256m",
		"--memory-swap=256m",
		"--cpus=0.5",
		"--network=none",
		"-e", "CODE="+code,
		image,
	)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	result := Result{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	// timeout check
	if ctx.Err() == context.DeadlineExceeded {
		result.Error = "Time Limit Exceeded (1s)"
		return result
	}

	if err != nil {
		result.Error = err.Error()
	}

	return result
}
