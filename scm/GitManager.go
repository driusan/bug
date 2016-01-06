package scm

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"os"
	"os/exec"
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

func (a GitManager) Commit(dir bugs.Directory, commitMsg string) error {
	cmd := exec.Command("git", "add", "-A", string(dir))
	if err := cmd.Run(); err != nil {
        fmt.Printf("Could not add issues to be commited: %s?\n", err.Error())
		return err

	}
	cmd = exec.Command("git", "commit", "-o", string(dir), "-m", commitMsg, "-q")
	if err := cmd.Run(); err != nil {
		// If nothing was added commit will have an error,
		// but we don't care it just means there's nothing
		// to commit.
		fmt.Printf("No new issues commited\n")
		return nil
	}
	return nil
}

func (a GitManager) GetSCMType() string {
	return "git"
}
