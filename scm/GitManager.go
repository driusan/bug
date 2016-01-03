package scm

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"os"
	"os/exec"
	"strings"
)

type PreconditionFailed string

func (a PreconditionFailed) Error() string {
	return string(a)
}

type ExecutionFailed string

func (a ExecutionFailed) Error() string {
	return string(a)
}

type GitManager struct{}

func (a GitManager) Purge(dir bugs.Directory) error {
	cmd := exec.Command("git", "clean", "-fd", string(dir))

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		return err
	}
	return nil
}

// This will try and commit the directory transparently,
// without affecting the current index.
// It does the following steps (after some validation):
//
// 1. "git stash create"
// 2. "git reset --mixed" (unstage the user's currently staged files)
// 3. "git add $(bug pwd)"
// 4. "git commit"
// 5a. "git reset --hard" (if there was any stash created,
// 						this is necessary for 5b to work.)
// 5b. "git stash apply --index" the stash from step 1
func (a GitManager) Commit(dir bugs.Directory, commitMsg string) error {
	type FileStatus struct {
		IndexStatus      string
		WorkingDirStatus string
		Filename         string
	}
	statusOutput := func(dir bugs.Directory) []FileStatus {
		cmd := exec.Command("git", "status", "--porcelain", "-z", string(dir))
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("Could not check git status")
			return nil
		}
		fileStatusLines := strings.Split(string(output), "\000")
		var files []FileStatus
		for _, line := range fileStatusLines {
			if len(line) == 0 {
				continue
			}
			files = append(files, FileStatus{
				IndexStatus:      line[0:1],
				WorkingDirStatus: line[1:2],
				Filename:         line[2:],
			})
		}
		return files
	}

	// Before doing anything, check git status to see if
	// the index is in a state that's going to cause an
	// error
	sOutput := statusOutput(bugs.GetIssuesDir())
	for _, file := range sOutput {
		if file.IndexStatus == "D" {
			return PreconditionFailed(fmt.Sprintf("You have manually staged changes in your issue directory which will conflict with %s commit.\n", os.Args[0]))
		}
	}

	sOutput = statusOutput(bugs.GetRootDir())
	for _, file := range sOutput {
		if file.IndexStatus == "A" {
			return PreconditionFailed(fmt.Sprintf("You have a new file staged in your git index, which will cause conflicts with %s commit. Please either commit your changes or unstage %s.\n", os.Args[0], file.Filename))
		}
	}

	cmd := exec.Command("git", "stash", "create")

	output, err := cmd.Output()

	if err != nil {
		return ExecutionFailed("Could not execute git stash create")
	}
	var stashHash string = strings.Trim(string(output), "\n")

	// Unstage everything, if there was anything stashed, so that
	// we don't commit things that the user has staged that aren't
	// issues
	if stashHash != "" {
		cmd = exec.Command("git", "reset", "--mixed")
		err = cmd.Run()

		if err != nil {
		}
	}

	// Commit the issues directory
	// git add $(bug pwd)
	// git commit -m "Added new issues" -q
	cmd = exec.Command("git", "add", "-A", string(dir))
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Could not add to index?\n")
	}
	cmd = exec.Command("git", "commit", "-m", commitMsg, "-q")
	err = cmd.Run()
	if err != nil {
		// If nothing was added commit will have an error,
		// but we don't care it just means there's nothing
		// to commit.
		fmt.Printf("No new issues commited\n")
	}

	// There were changes that had been stashed, so we need
	// to restore them with git stash apply.. first, we
	// need to do a "git reset --hard" so that the dirty working
	// tree doesn't cause an error. This isn't as scary as it
	// sounds, since immediately after git reset --hard we apply
	// a stash which has the exact same changes that we just threw
	// away.
	if stashHash != "" {
		cmd = exec.Command("git", "reset", "--hard")
		err = cmd.Run()
		if err != nil {
			return ExecutionFailed(fmt.Sprintf("Error resetting the git working tree\nThe stash which should have your changes is: %s\n", stashHash))
		}
		cmd = exec.Command("git", "stash", "apply", "--index", stashHash)
		err = cmd.Run()
		if err != nil {
			return ExecutionFailed(fmt.Sprintf("Error restoring the git working tree\nThe stash which should have your changes is: %s\n", stashHash))
		}
	}
	return nil
}

func (a GitManager) GetSCMType() string {
	return "git"
}
