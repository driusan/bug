package scm

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type HgCommit struct {
	commit string
	log    string
}

func (c HgCommit) CommitID() string {
	return c.commit
}
func (c HgCommit) LogMsg() string {
	return c.log
}
func (c HgCommit) Diff() (string, error) {
	return runCmd("hg", "log", "-p", "-g", "-r", c.commit, "--template={changelog}")
}

type HgTester struct {
	handler SCMHandler
	workdir string
}

func (t HgTester) GetLogs() ([]Commit, error) {
	logs, err := runCmd("hg", "log", "-r", ":", "--template", "{node} {desc}\\n")
	if err != nil {
		return nil, err
	}
	logMsgs := strings.Split(logs, "\n")
	// the last line is empty, so don't allocate 1 for
	// it
	commits := make([]Commit, len(logMsgs)-1)
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

func (c HgTester) AssertStagingIndex(t *testing.T, f []FileStatus) {
	for _, file := range f {
		out, err := runCmd("hg", "status", file.Filename)
		if err != nil {
			t.Error("Could not get status of " + file.Filename)
		}

		// hg status doesn't include the working directory status
		expected := file.IndexStatus + " " + file.Filename + "\n"
		if out != expected {
			t.Error("Unexpected status. Got " + out + " not " + expected)
		}

	}
}

func (c HgTester) StageFile(file string) error {
	_, err := runCmd("hg", "add", file)
	return err
}
func (c *HgTester) Setup() error {
	if dir, err := ioutil.TempDir("", "hgbug"); err == nil {
		c.workdir = dir
		os.Chdir(c.workdir)
	} else {
		return err
	}

	_, err := runCmd("hg", "init")
	if err != nil {
		return err
	}
	c.handler = HgManager{}
	return nil
}

func (c HgTester) TearDown() {
	os.RemoveAll(c.workdir)
}

func (c HgTester) GetWorkDir() string {
	return c.workdir
}
func (c HgTester) AssertCleanTree(t *testing.T) {
	out, err := runCmd("hg", "status")
	if err != nil {
		t.Error("Error running hg status")
	}
	if out != "" {
		fmt.Printf("\"%s\"\n", out)
		t.Error("Unexpected Output from hg status")
	}
}

func (m HgTester) GetManager() SCMHandler {
	return m.handler
}

func TestHgBugRenameCommits(t *testing.T) {
	tester := HgTester{}

	expectedDiffs := []string{
		`diff --git a/issues/Test-bug/Description b/issues/Test-bug/Description
new file mode 100644

`, `diff --git a/issues/Renamed-bug/Description b/issues/Renamed-bug/Description
new file mode 100644
diff --git a/issues/Test-bug/Description b/issues/Test-bug/Description
deleted file mode 100644

`}
	runtestRenameCommitsHelper(&tester, t, expectedDiffs)
}
func TestHgFilesOutsideOfBugNotCommited(t *testing.T) {
	tester := HgTester{}
	tester.handler = HgManager{}
	runtestCommitDirtyTree(&tester, t)
}

func TestHgGetType(t *testing.T) {
	m := HgManager{}

	if m.GetSCMType() != "hg" {
		t.Error("Incorrect type for HgManager")
	}
}

func TestHgPurge(t *testing.T) {
	// This should eventually be replaced by something more
	// like:
	//      m := HgTester{}
	//      m.handler = HgManager{}
	//      runtestPurgeFiles(&m, t)
	// but since the current behaviour is to return a not
	// supported error, that would evidently fail..
	m := HgManager{}
	err := m.Purge("/tmp/imaginaryHgRepo")

	switch err.(type) {
	case UnsupportedType:
		// This is valid, do nothing.
	default:
		t.Error("Unexpected return value for Hg purge function.")
	}

}
