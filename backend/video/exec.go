package video

import (
	"bytes"
	"fmt"
	"os/exec"
)

func _exec(name string, args ...string) ([]byte, error) {

	cmd := exec.Command(name, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("error: %s; stdout: %s, stderr: %s", err, stdout.String(), stderr.String())
	}

	return stdout.Bytes(), nil
}
