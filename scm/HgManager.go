package scm

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"os/exec"
)

type HgManager struct{}

func (a HgManager) Purge(dir bugs.Directory) error {
	return UnsupportedType("Purge is not supported under Hg. Sorry!")
}

func (a HgManager) Commit(dir bugs.Directory, commitMsg string) error {
	cmd := exec.Command("hg", "addremove", string(dir))
	if err := cmd.Run(); err != nil {
		fmt.Printf("Could not add issues to be commited: %s?\n", err.Error())
		return err
	}

	cmd = exec.Command("hg", "commit", string(dir), "-m", commitMsg)
	if err := cmd.Run(); err != nil {
		fmt.Printf("No new issues to commit.\n")
		return err
	}
	return nil
}

func (a HgManager) GetSCMType() string {
	return "hg"
}
