package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// exit handles the common error exit case
func exit(err error) {
	if exitErr, ok := err.(*exec.ExitError); ok {
		fmt.Println(exitErr.Stderr)
	}

	os.Exit(1)
}

// head returns the current head for the repo
func head() string {
	headRef, headErr := exec.Command(
		"git",
		"rev-parse",
		"--abbrev-ref",
		"HEAD",
	).Output()

	if headErr != nil {
		exit(headErr)
	}

	return strings.TrimSpace(string(headRef))
}

// status returns the current stats of the repo
func status() string {
	statusOut, statusErr := exec.Command(
		"git",
		"status",
		"-s",
	).Output()

	if statusErr != nil {
		exit(statusErr)
	}

	return string(statusOut)
}

// changes returns the number of lines in the status output
func changes() int {
	statusOut := status()

	changes := strings.Split(statusOut, "\n")
	count := 0

	for _, line := range changes {
		line = strings.TrimSpace(line)
		if line != "" {
			count++
		}
	}

	return count
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: git pr [name] [rev]")
		os.Exit(1)
	}

	name := os.Args[1]
	rev := os.Args[2]

	headRev := head()
	changeCount := changes()

	if changeCount > 0 {
		fmt.Println("Cannot continue: You have pending changes")
		os.Exit(1)
	}

	// create the new branch
	branchCmd := exec.Command(
		"git",
		"branch",
		name,
		"origin/master",
		"--no-track",
	)
	branchCmd.Stdout = os.Stdout
	branchCmd.Stderr = os.Stderr

	branchErr := branchCmd.Run()
	if branchErr != nil {
		os.Exit(1)
	}

	// check out the new branch
	checkoutCmd := exec.Command(
		"git",
		"checkout",
		"-q",
		name,
	)
	checkoutCmd.Stdout = os.Stdout
	checkoutCmd.Stderr = os.Stderr

	checkoutErr := checkoutCmd.Run()
	if checkoutErr != nil {
		os.Exit(1)
	}

	// cherry pick
	cherrypickCmd := exec.Command(
		"git",
		"cherry-pick",
		rev,
	)
	cherrypickCmd.Stdout = os.Stdout
	cherrypickCmd.Stderr = os.Stderr

	cherrypickErr := cherrypickCmd.Run()
	if cherrypickErr != nil {
		os.Exit(1)
	}

	// push
	pushCmd := exec.Command(
		"git",
		"push",
		"origin",
		name,
	)
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr

	pushErr := pushCmd.Run()
	if pushErr != nil {
		os.Exit(1)
	}

	// check out original head
	checkoutOrigCmd := exec.Command(
		"git",
		"checkout",
		"-q",
		headRev,
	)
	checkoutOrigCmd.Stdout = os.Stdout
	checkoutOrigCmd.Stderr = os.Stderr

	checkoutOrigErr := checkoutOrigCmd.Run()
	if checkoutOrigErr != nil {
		os.Exit(1)
	}

	// delete the branch that was created
	deleteBranchCmd := exec.Command(
		"git",
		"branch",
		"-D",
		name,
	)
	deleteBranchCmd.Stdout = os.Stdout
	deleteBranchCmd.Stderr = os.Stderr

	deleteBranchErr := deleteBranchCmd.Run()
	if deleteBranchErr != nil {
		os.Exit(1)
	}
}
