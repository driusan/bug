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

    ioutil.WriteFile(dir + "/issues/Test/Identifier", []byte("TestBug\n"), 0600)

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
