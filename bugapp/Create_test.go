package bugapp

import (
	"fmt"
	//	"io"
	"io/ioutil"
	"os"
	"testing"
)

func captureOutput(f func(), t *testing.T) (string, string) {
	// Capture STDOUT with a pipe
	stdout := os.Stdout
	stderr := os.Stderr
	so, op, _ := os.Pipe() //outpipe
	oe, ep, _ := os.Pipe() //errpipe
	defer func(stdout, stderr *os.File) {
		os.Stdout = stdout
		os.Stderr = stderr
	}(stdout, stderr)

	os.Stdout = op
	os.Stderr = ep

	f()

	os.Stdout = stdout
	os.Stderr = stderr

	op.Close()
	ep.Close()

	errOutput, err := ioutil.ReadAll(oe)
	if err != nil {
		t.Error("Could not get output from stderr")
	}
	stdOutput, err := ioutil.ReadAll(so)
	if err != nil {
		t.Error("Could not get output from stdout")
	}
	return string(stdOutput), string(errOutput)
}

// Captures stdout and stderr to ensure that
// a usage line gets printed to Stderr when
// no parameters are specified
func TestCreateHelpOutput(t *testing.T) {

	stdout, stderr := captureOutput(func() {
		Create(ArgumentList{})
	}, t)

	if stdout != "" {
		t.Error("Unexpected output on stdout.")
	}
	if stderr[:7] != "Usage: " {
		t.Error("Expected usage information with no arguments")
	}

}

// Test a very basic invocation of "Create" with the -n
// argument. We can't try it without -n, since it means
// an editor will be spawned..
func TestCreateNoEditor(t *testing.T) {
	dir, err := ioutil.TempDir("", "createtest")
	if err != nil {
		t.Error("Could not create temporary dir for test")
		return
	}
	os.Chdir(dir)
	os.MkdirAll("issues", 0700)
	defer os.RemoveAll(dir)

	stdout, stderr := captureOutput(func() {
		Create(ArgumentList{"-n", "Test", "bug"})
	}, t)
	if stderr != "" {
		t.Error("Unexpected output on STDERR for Test-bug")
	}
	if stdout != "Created issue: Test bug\n" {
		t.Error("Unexpected output on STDOUT for Test-bug")
	}
	issuesDir, err := ioutil.ReadDir(fmt.Sprintf("%s/issues/", dir))
	if err != nil {
		t.Error("Could not read issues directory")
		return
	}
	if len(issuesDir) != 1 {
		t.Error("Unexpected number of issues in issues dir\n")
	}

	bugDir, err := ioutil.ReadDir(fmt.Sprintf("%s/issues/Test-bug", dir))
	if len(bugDir) != 1 {
		t.Error("Unexpected number of files found in Test-bug dir\n")
	}
	if err != nil {
		t.Error("Could not read Test-bug directory")
		return
	}

	file, err := ioutil.ReadFile(fmt.Sprintf("%s/issues/Test-bug/Description", dir))
	if err != nil {
		t.Error("Could not load description file for Test bug" + err.Error())
	}
	if len(file) != 0 {
		t.Error("Expected empty file for Test bug")
	}
}
