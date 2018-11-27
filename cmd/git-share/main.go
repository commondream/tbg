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

// run runs the given command with stdout and stderr mapped to this process
func run(program string, args ...string) error {
	cmd := exec.Command(program, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
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

// isRef returns whether or not the ref exists
func isRef(ref string) bool {
	refErr := exec.Command(
		"git",
		"show-ref",
		ref,
	).Run()

	if refErr != nil {
		return false
	}

	return true
}

func revParse(rev string) string {
	revParseOut, revParseErr := exec.Command(
		"git",
		"rev-parse",
		rev,
	).Output()

	if revParseErr != nil {
		exit(revParseErr)
	}

	return string(strings.TrimSpace(string(revParseOut)))
}

func main() {
	fmt.Println("...")

	if len(os.Args) < 3 {
		fmt.Println("Usage: git share [name] [rev]")
		os.Exit(1)
	}

	name := os.Args[1]
	rev := os.Args[2]

	headRev := head()
	changeCount := changes()

	if changeCount > 0 {
		fmt.Println("Cannot continue: pending changes")
		os.Exit(1)
	}

	// rev-parse the rev param so things like HEAD are possible
	revParsed := revParse(rev)

	if !isRef(fmt.Sprintf("origin/%s", name)) {
		// create the new branch
		fmt.Printf("* Creating share branch %s\n", name)
		branchErr := run(
			"git",
			"branch",
			name,
			"origin/master",
			"--no-track",
		)
		if branchErr != nil {
			os.Exit(1)
		}
		fmt.Println()
	}

	// check out the new branch
	fmt.Printf("*  Checking out %s\n", name)
	checkoutErr := run(
		"git",
		"checkout",
		"-q",
		name,
	)
	if checkoutErr != nil {
		os.Exit(1)
	}
	fmt.Println()

	// cherry pick
	fmt.Printf("* Picking %s onto %s\n", revParsed, name)
	cherrypickErr := run(
		"git",
		"cherry-pick",
		revParsed,
	)
	if cherrypickErr != nil {
		os.Exit(1)
	}
	fmt.Println()

	// push
	fmt.Printf("* Pushing %s\n", name)
	pushErr := run(
		"git",
		"push",
		"origin",
		name,
	)
	if pushErr != nil {
		os.Exit(1)
	}
	fmt.Println()

	// check out original head
	fmt.Printf("* Back to %s\n", headRev)
	checkoutOrigErr := run(
		"git",
		"checkout",
		"-q",
		headRev,
	)
	if checkoutOrigErr != nil {
		os.Exit(1)
	}
	fmt.Println()

	// delete the branch that was created
	fmt.Printf("* Deleting %s\n", name)
	deleteBranchErr := run(
		"git",
		"branch",
		"-D",
		name,
	)
	if deleteBranchErr != nil {
		os.Exit(1)
	}
	fmt.Println()
}
