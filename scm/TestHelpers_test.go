package scm

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

type Commit interface {
	CommitID() string
	LogMsg() string
	Diff() (string, error)
}

type ManagerTester interface {
	GetLogs() ([]Commit, error)
	Setup() error
	GetWorkDir() string
	TearDown()
	StageFile(string) error
	AssertCleanTree(t *testing.T)
	AssertStagingIndex(*testing.T, []FileStatus)
	GetManager() SCMHandler
}

func runCmd(cmd string, options ...string) (string, error) {
	runcmd := exec.Command(cmd, options...)
	out, err := runcmd.CombinedOutput()

	return string(out), err
}

func assertLogs(tester ManagerTester, t *testing.T, titles []map[string]bool, diffs []string) {
	logs, err := tester.GetLogs()
	if err != nil {
		t.Error("Could not get scm logs" + err.Error())
		return
	}

	if len(diffs) != len(titles) {
		t.Error("Different number of diffs from titles")
		return
	}
	if len(logs) != len(titles) || len(logs) != len(diffs) {
		t.Error("Unexpected number of log messages")
		return
	}

	for i, _ := range titles {
		if _, ok := titles[i][logs[i].LogMsg()]; !ok {
			t.Error("Unexpected commit message:" + logs[i].LogMsg())
		}

		if diff, err := logs[i].Diff(); err != nil {
			t.Error("Could not get diff of commit")
		} else {
			if diff != diffs[i] {
				// get shortest commit msg to keep errors simple
				var s string
				for k := range titles[i] {
					if len(s) == 0 || len(k) < len(s) {
						s = k
					}
				}
				t.Error(fmt.Sprintf("Incorrect diff for i=%d, title=%s", i, s))
				fmt.Fprintf(os.Stderr, "Got %s, expected %s", diff, diffs[i])
			}
		}
	}
}

func runtestRenameCommitsHelper(tester ManagerTester, t *testing.T, expectedDiffs []string) {
	err := tester.Setup()
	defer tester.TearDown()
	if err != nil {
		t.Error("Could not initialize repo")
		return
	}

	m := tester.GetManager()
	if m == nil {
		t.Error("Could not get manager")
		return
	}
	os.MkdirAll("issues/Test-bug", 0755)
	ioutil.WriteFile("issues/Test-bug/Description", []byte(""), 0644)
	m.Commit(bugs.Directory(tester.GetWorkDir()), "Initial commit")
	runCmd("bug", "relabel", "1", "Renamed", "bug")
	m.Commit(bugs.Directory(tester.GetWorkDir()), "This is a test rename")

	tester.AssertCleanTree(t)

	assertLogs(tester, t, []map[string]bool{{
		"Initial commit":          true, // simple format
		`Create issue "Test-bug"`: true, // rich format
	}, {
		"This is a test rename":                    true, // simple format
		`Update issues: "Test-bug", "Renamed-bug"`: true, // rich format
		`Update issues: "Renamed-bug", "Test-bug"`: true, // has two alternatives equally good
	}}, expectedDiffs)

}
func runtestCommitDirtyTree(tester ManagerTester, t *testing.T) {
	err := tester.Setup()
	if err != nil {
		panic("Something went wrong trying to initialize git:" + err.Error())
	}
	defer tester.TearDown()
	m := tester.GetManager()
	if m == nil {
		t.Error("Could not get manager")
		return
	}
	os.Mkdir("issues", 0755)
	runCmd("bug", "create", "-n", "Test", "bug")
	if err = ioutil.WriteFile("donotcommit.txt", []byte(""), 0644); err != nil {
		t.Error("Could not write file")
		return
	}
	tester.AssertStagingIndex(t, []FileStatus{
		FileStatus{"donotcommit.txt", "?", "?"},
	})

	m.Commit(bugs.Directory(tester.GetWorkDir()+"/issues"), "Initial commit")
	tester.AssertStagingIndex(t, []FileStatus{
		FileStatus{"donotcommit.txt", "?", "?"},
	})
	tester.StageFile("donotcommit.txt")
	tester.AssertStagingIndex(t, []FileStatus{
		FileStatus{"donotcommit.txt", "A", " "},
	})
	m.Commit(bugs.Directory(tester.GetWorkDir()+"/issues"), "Initial commit")
	tester.AssertStagingIndex(t, []FileStatus{
		FileStatus{"donotcommit.txt", "A", " "},
	})
}

func runtestPurgeFiles(tester ManagerTester, t *testing.T) {
	err := tester.Setup()
	if err != nil {
		panic("Something went wrong trying to initialize: " + err.Error())
	}
	defer tester.TearDown()
	m := tester.GetManager()
	if m == nil {
		t.Error("Could not get manager")
		return
	}
	os.Mkdir("issues", 0755)
	// Commit a bug which should stay around after the purge
	runCmd("bug", "create", "-n", "Test", "bug")
	m.Commit(bugs.Directory(tester.GetWorkDir()+"/issues"), "Initial commit")

	// Create another bug to elimate with "bug purge"
	runCmd("bug", "create", "-n", "Test", "purge", "bug")
	err = m.Purge(bugs.Directory(tester.GetWorkDir() + "/issues"))
	if err != nil {
		t.Error("Error purging bug directory: " + err.Error())
	}
	issuesDir, err := ioutil.ReadDir("issues") //fmt.Sprintf("%s/issues/", tester.GetWorkDir()))
	if err != nil {
		t.Error("Error reading issues directory")
	}
	if len(issuesDir) != 1 {
		t.Error("Unexpected number of directories in issues/ after purge.")
	}
	if len(issuesDir) > 0 && issuesDir[0].Name() != "Test-bug" {
		t.Error("Expected Test-bug to remain.")
	}
}
