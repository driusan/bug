package scm

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var workdir string
var handler SCMHandler

type GitCommit struct {
	CommitID string
	LogMsg   string
}

func (c GitCommit) Diff() (string, error) {
	return runCmd("git", "show", "--pretty=format:%b", c.CommitID)
}
func getLogs() ([]GitCommit, error) {
	logs, err := runCmd("git", "log", "--pretty=oneline", "--reverse")
	if err != nil {
		return nil, err
	}
	logMsgs := strings.Split(logs, "\n")
	// the last line is empty, so don't allocate 1 for
	// it
	commits := make([]GitCommit, len(logMsgs)-1)
	for idx, commitText := range logMsgs {
		if commitText == "" {
			continue
		}
		spaceIdx := strings.Index(commitText, " ")
		if spaceIdx >= 0 {
			commits[idx] = GitCommit{commitText[0:spaceIdx], commitText[spaceIdx+1:]}
		}
	}
	return commits, nil
}
func runCmd(cmd string, options ...string) (string, error) {
	runcmd := exec.Command(cmd, options...)
	out, err := runcmd.CombinedOutput()

	return string(out), err
}

func setupGit() error {
	if dir, err := ioutil.TempDir("", "gitbug"); err == nil {
		workdir = dir
		os.Chdir(workdir)
	} else {
		return err
	}

	_, err := runCmd("git", "init")
	if err != nil {
		return err
	}

	handler = GitManager{}
	return nil
}
func tearDownGit() {
	os.RemoveAll(workdir)
}

func assertCleanGitTree(t *testing.T) {
	out, err := runCmd("git", "status", "--porcelain")
	if err != nil {
		t.Error("Error running git status")
	}
	if out != "" {
		t.Error("Unexpected Output from git status")
	}
}

func TestGitBugRenameCommits(t *testing.T) {
	err := setupGit()
	if err != nil {
		panic("Something went wrong trying to initialize git:" + err.Error())
	}
	defer tearDownGit()

	os.Mkdir("issues", 0755)
	runCmd("bug", "create", "-n", "Test bug")
	handler.Commit(bugs.Directory(workdir), "Initial commit")
	runCmd("bug", "relabel", "1", "Renamed bug")
	handler.Commit(bugs.Directory(workdir), "This is a test rename")

	assertCleanGitTree(t)

	logs, err := getLogs()
	if err != nil {
		t.Error("Could not get git logs")
	}

	if len(logs) != 2 {
		fmt.Printf("Got %d log messages. %s\n", len(logs), logs)
		t.Error("Unexpected number of log messages")
	}

	if logs[0].LogMsg != "Initial commit" {
		t.Error("Unexpected commit message:" + logs[0].LogMsg)
	}
	if logs[1].LogMsg != "This is a test rename" {
		t.Error("Unexpected commit message:" + logs[1].LogMsg)
	}

	diff, err := logs[0].Diff()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		t.Error("Could not get diff of commit")
	}
	if diff != `
diff --git a/issues/Test-bug/Description b/issues/Test-bug/Description
new file mode 100644
index 0000000..e69de29
` {
		fmt.Printf("Got: \"%s\"\n", diff)
		t.Error("Incorrect diff")
	}

	diff, err = logs[1].Diff()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		t.Error("Could not get diff of commit")
	}
	if diff != `
diff --git a/issues/Renamed-bug/Description b/issues/Renamed-bug/Description
new file mode 100644
index 0000000..e69de29
diff --git a/issues/Test-bug/Description b/issues/Test-bug/Description
deleted file mode 100644
index e69de29..0000000
` {
		fmt.Printf("Got: \"%s\"\n", diff)
		t.Error("Incorrect diff")
	}
}
