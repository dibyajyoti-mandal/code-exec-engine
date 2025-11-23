package executor

import (
	"bytes"
	"os/exec"
)

type Result struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Error  string `json:"error"`
}

func RunInDocker(image string, code string) Result {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command(
		"docker", "run", "--rm",
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

	if err != nil {
		result.Error = err.Error()
	}

	return result
}
