package main

import (
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command(
		"git",
		"log",
		"master",
		"^origin/master",
		"--no-merges",
		"--pretty=oneline",
		"--abbrev-commit",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		os.Exit(1)
	}
}
