package scm

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type GitCommit struct {
	commit string
	log    string
}

func (c GitCommit) CommitID() string {
	return c.commit
}
func (c GitCommit) LogMsg() string {
	return c.log
}
func (c GitCommit) Diff() (string, error) {
	return runCmd("git", "show", "--pretty=format:%b", c.CommitID())
}

type GitTester struct {
	handler SCMHandler
	workdir string
}

func (t GitTester) GetLogs() ([]Commit, error) {
	logs, err := runCmd("git", "log", "--oneline", "--reverse", "-z")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving git logs: %s", logs)
		return nil, err
	}
	logMsgs := strings.Split(logs, "\000")
	// the last line is empty, so don't allocate 1 for
	// it
	commits := make([]Commit, len(logMsgs)-1)
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

func (g GitTester) AssertStagingIndex(t *testing.T, f []FileStatus) {
	for _, file := range f {
		out, err := runCmd("git", "status", "--porcelain", file.Filename)
		if err != nil {
			t.Error("Could not run git status")
		}
		expected := file.IndexStatus + file.WorkingStatus + " " + file.Filename + "\n"
		if out != expected {
			t.Error("Incorrect file status")
			t.Error("Got" + out + " not " + expected)
		}
	}
}

func (g GitTester) StageFile(file string) error {
	_, err := runCmd("git", "add", file)
	return err
}
func (t *GitTester) Setup() error {
	if dir, err := ioutil.TempDir("", "gitbug"); err == nil {
		t.workdir = dir
		os.Chdir(t.workdir)
	} else {
		return err
	}

	out, err := runCmd("git", "init", ".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing git: %s", out)
		return err
	}

	t.handler = GitManager{}
	return nil
}
func (t GitTester) TearDown() {
	os.RemoveAll(t.workdir)
}
func (t GitTester) GetWorkDir() string {
	return t.workdir
}

func (m GitTester) AssertCleanTree(t *testing.T) {
	out, err := runCmd("git", "status", "--porcelain")
	if err != nil {
		t.Error("Error running git status")
	}
	if out != "" {
		t.Error("Unexpected Output from git status (expected nothing):\n" + out)
	}
}

func (m GitTester) GetManager() SCMHandler {
	return m.handler
}

func TestGitBugRenameCommits(t *testing.T) {
	if os.Getenv("TRAVIS") == "true" {
		t.Skip("Skipping test which fails only under Travis for unknown reasons..")
		return
	}
	gm := GitTester{}
	gm.handler = GitManager{}

	expectedDiffs := []string{
		`
diff --git a/issues/Test-bug/Description b/issues/Test-bug/Description
new file mode 100644
index 0000000..e69de29
`, `
diff --git a/issues/Renamed-bug/Description b/issues/Renamed-bug/Description
new file mode 100644
index 0000000..e69de29
diff --git a/issues/Test-bug/Description b/issues/Test-bug/Description
deleted file mode 100644
index e69de29..0000000
`}

	runtestRenameCommitsHelper(&gm, t, expectedDiffs)
}

func TestGitFilesOutsideOfBugNotCommited(t *testing.T) {
	gm := GitTester{}
	gm.handler = GitManager{}
	runtestCommitDirtyTree(&gm, t)
}
