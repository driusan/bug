package bugapp

import (
	"fmt"
	//	"io"
	"io/ioutil"
	"os"
	"testing"
)

// Captures stdout and stderr to ensure that
// a usage line gets printed to Stderr when
// no parameters are specified
func TestCloseHelpOutput(t *testing.T) {

	stdout, stderr := captureOutput(func() {
		Close(ArgumentList{})
	}, t)

	if stdout != "" {
		t.Error("Unexpected output on stdout.")
	}
	if stderr[:7] != "Usage: " {
		t.Error("Expected usage information with no arguments")
	}

}

// Test closing a bug given it's directory index
func TestCloseByIndex(t *testing.T) {
	dir, err := ioutil.TempDir("", "closetest")
	defer os.RemoveAll(dir)
	if err != nil {
		t.Error("Could not create temporary dir for test")
		return
	}
	os.Chdir(dir)
	os.MkdirAll("issues/Test", 0700)

	// On MacOS, /tmp is a symlink, which causes GetDirectory() to return
	// a different path than expected in these tests, so make the issues
	// directory explicit with an environment variable
	err = os.Setenv("PMIT", dir)
	if err != nil {
		t.Error("Could not set environment variable: " + err.Error())
		return
	}

	ioutil.WriteFile(dir+"/issues/Test/Identifier", []byte("TestBug\n"), 0600)

	issuesDir, err := ioutil.ReadDir(fmt.Sprintf("%s/issues/", dir))
	// Assert that there's 1 bug to start, otherwise what are we closing?
	if err != nil || len(issuesDir) != 1 {
		t.Error("Could not read issues directory")
		return
	}
	stdout, stderr := captureOutput(func() {
		Close(ArgumentList{"TestBug"})
	}, t)
	if stderr != "" {
		t.Error("Unexpected output on STDERR for Test-bug")
	}
	if stdout != fmt.Sprintf("Removing %s/issues/Test\n", dir) {
		t.Error("Unexpected output on STDOUT for Test-bug")
		fmt.Printf("Got %s\nExpected%s", stdout, fmt.Sprintf("Removing %s/issues/Test\n", dir))
	}
	issuesDir, err = ioutil.ReadDir(fmt.Sprintf("%s/issues/", dir))
	if err != nil {
		t.Error("Could not read issues directory")
		return
	}
	// After closing, there should be 0 bugs.
	if len(issuesDir) != 0 {
		t.Error("Unexpected number of issues in issues dir\n")
	}
}

func TestCloseBugByIdentifier(t *testing.T) {
	dir, err := ioutil.TempDir("", "close")
	if err != nil {
		t.Error("Could not create temporary dir for test")
		return
	}
	os.Chdir(dir)
	os.MkdirAll("issues/Test", 0700)
	defer os.RemoveAll(dir)

	// On MacOS, /tmp is a symlink, which causes GetDirectory() to return
	// a different path than expected in these tests, so make the issues
	// directory explicit with an environment variable
	err = os.Setenv("PMIT", dir)
	if err != nil {
		t.Error("Could not set environment variable: " + err.Error())
		return
	}

	issuesDir, err := ioutil.ReadDir(fmt.Sprintf("%s/issues/", dir))
	// Assert that there's 1 bug to start, otherwise what are we closing?
	if err != nil || len(issuesDir) != 1 {
		t.Error("Could not read issues directory")
		return
	}
	stdout, stderr := captureOutput(func() {
		Close(ArgumentList{"1"})
	}, t)
	if stderr != "" {
		t.Error("Unexpected output on STDERR for Test-bug")
	}
	if stdout != fmt.Sprintf("Removing %s/issues/Test\n", dir) {
		t.Error("Unexpected output on STDOUT for Test-bug")
		fmt.Printf("Got %s\nExpected: %s\n", stdout, dir)
	}
	issuesDir, err = ioutil.ReadDir(fmt.Sprintf("%s/issues/", dir))
	if err != nil {
		t.Error("Could not read issues directory")
		return
	}
	// After closing, there should be 0 bugs.
	if len(issuesDir) != 0 {
		t.Error("Unexpected number of issues in issues dir\n")
	}
}

func TestCloseMultipleIndexesWithLastIndex(t *testing.T) {
	dir, err := ioutil.TempDir("", "closetest")
	defer os.RemoveAll(dir)
	if err != nil {
		t.Error("Could not create temporary dir for test")
		return
	}
	os.Chdir(dir)
	os.Setenv("PMIT", dir)
	os.MkdirAll("issues/Test", 0700)
	os.MkdirAll("issues/Test2", 0700)
	os.MkdirAll("issues/Test3", 0700)
	issuesDir, err := ioutil.ReadDir(fmt.Sprintf("%s/issues/", dir))
	if err != nil {
		t.Error("Could not read issues directory")
		return
	}
	if len(issuesDir) != 3 {
		t.Error("Unexpected number of issues in issues dir after creating multiple issues\n")
	}
	_, stderr := captureOutput(func() {
		Close(ArgumentList{"1", "3"})
	}, t)
	issuesDir, err = ioutil.ReadDir(fmt.Sprintf("%s/issues/", dir))
	if err != nil {
		t.Error("Could not read issues directory")
		return
	}
	// After closing, there should be 1 bug. Otherwise, it probably
	// means that the last error was "invalid index" since indexes
	// were renumbered after closing the first bug.
	if len(issuesDir) != 1 {
		fmt.Printf("%s\n\n", stderr)
		t.Error("Unexpected number of issues in issues dir after closing multiple issues\n")
	}
}

func TestCloseMultipleIndexesAtOnce(t *testing.T) {
	dir, err := ioutil.TempDir("", "closetest")
	defer os.RemoveAll(dir)
	if err != nil {
		t.Error("Could not create temporary dir for test")
		return
	}
	os.Chdir(dir)
	os.Setenv("PMIT", dir)
	os.MkdirAll("issues/Test", 0700)
	os.MkdirAll("issues/Test2", 0700)
	os.MkdirAll("issues/Test3", 0700)
	issuesDir, err := ioutil.ReadDir(fmt.Sprintf("%s/issues/", dir))
	if err != nil {
		t.Error("Could not read issues directory")
		return
	}
	if len(issuesDir) != 3 {
		t.Error("Unexpected number of issues in issues dir after creating multiple issues\n")
	}
	_, _ = captureOutput(func() {
		Close(ArgumentList{"1", "2"})
	}, t)
	issuesDir, err = ioutil.ReadDir(fmt.Sprintf("%s/issues/", dir))
	if err != nil {
		t.Error("Could not read issues directory")
		return
	}
	if len(issuesDir) != 1 {
		t.Error("Unexpected number of issues in issues dir after closing multiple issues\n")
		return
	}

	// 1 and 2 should have closed. If 3 was renumbered after 1 was closed,
	// it would be closed instead.
	if issuesDir[0].Name() != "Test3" {
		t.Error("Closed incorrect issue when closing multiple issues.")
	}
}
