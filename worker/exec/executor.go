package executor

import (
	"bytes"
	"context"
	"os/exec"
	"time"

	"github.com/dibyajyoti-mandal/code-exec-engine/constants"
	"github.com/dibyajyoti-mandal/code-exec-engine/models"
)

func RunInDocker(image string, code string) models.Result {
	var stdout, stderr bytes.Buffer

	// timeout = 1 second
	ctx, cancel := context.WithTimeout(context.Background(), constants.CONTAINER_LIFE*time.Second)
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

	result := models.Result{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	// timeout check
	if ctx.Err() == context.DeadlineExceeded {
		result.Error = "Time Limit Exceeded"
		return result
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() == 124 {
			result.Error = "Time Limit Exceeded"
			return result
		}
	}

	return result
}
