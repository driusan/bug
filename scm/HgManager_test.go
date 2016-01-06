package scm

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type HgCommit struct {
	CommitID string
	LogMsg   string
}

func (c HgCommit) Diff() (string, error) {
	return runCmd("hg", "log", "-p", "-g", "-r", c.CommitID, "--template={changelog}")
}
func getHgLogs() ([]HgCommit, error) {
    logs, err := runCmd("hg", "log", "-r", ":", "--template", "{node} {desc}\\n")
	if err != nil {
		return nil, err
	}
	logMsgs := strings.Split(logs, "\n")
	// the last line is empty, so don't allocate 1 for
	// it
	commits := make([]HgCommit, len(logMsgs)-1)
	for idx, commitText := range logMsgs {
		if commitText == "" {
			continue
		}
		spaceIdx := strings.Index(commitText, " ")
		if spaceIdx >= 0 {
			commits[idx] = HgCommit{commitText[0:spaceIdx], commitText[spaceIdx+1:]}
		}
	}
	return commits, nil
}

func setupHg() error {
	if dir, err := ioutil.TempDir("", "hgbug"); err == nil {
		workdir = dir
		os.Chdir(workdir)
	} else {
		return err
	}

	_, err := runCmd("hg", "init")
	if err != nil {
		return err
	}

	handler = HgManager{}
	return nil
}
func tearDownHg() {
	os.RemoveAll(workdir)
}

func assertCleanHgTree(t *testing.T) {
	out, err := runCmd("hg", "status")
	if err != nil {
		t.Error("Error running hg status")
	}
	if out != "" {
        fmt.Printf("\"%s\"\n", out)
		t.Error("Unexpected Output from hg status")
	}
}

func TestHgBugRenameCommits(t *testing.T) {
	err := setupHg()
	if err != nil {
		panic("Something went wrong trying to initialize Hg:" + err.Error())
	}
	defer tearDownHg()


	os.Mkdir("issues", 0755)
	runCmd("bug", "create", "-n", "Test bug")
	handler.Commit(bugs.Directory(workdir), "Initial commit")
	runCmd("bug", "relabel", "1", "Renamed bug")
	handler.Commit(bugs.Directory(workdir), "This is a test rename")

	assertCleanHgTree(t)

	logs, err := getHgLogs()
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
	if diff != `diff --git a/issues/Test-bug/Description b/issues/Test-bug/Description
new file mode 100644

` {
		fmt.Printf("Got: \"%s\"\n", diff)
		t.Error("Incorrect diff")
	}

	diff, err = logs[1].Diff()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		t.Error("Could not get diff of commit")
	}
	if diff != `diff --git a/issues/Renamed-bug/Description b/issues/Renamed-bug/Description
new file mode 100644
diff --git a/issues/Test-bug/Description b/issues/Test-bug/Description
deleted file mode 100644

` {
		fmt.Printf("Got: \"%s\"\n", diff)
		t.Error("Incorrect diff")
	}
}
