package main

import (
	"bytes"
	"os/exec"
)

func runOpRead(opPath, itemURL string) (output, errOutput string, err error) {
	cmd := exec.Command(opPath, "read", "-n", itemURL)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	outputBytes, err := cmd.Output()
	return string(outputBytes),
		string(bytes.TrimRight(stderr.Bytes(), "\n")),
		err
}
