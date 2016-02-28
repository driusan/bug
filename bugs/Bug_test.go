package bugs

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

type tester struct {
	dir string
	bug *Bug
}

func (t *tester) Setup() {
	gdir, err := ioutil.TempDir("", "issuetest")
	if err == nil {
		os.Chdir(gdir)
		t.dir = gdir
		os.Unsetenv("PMIT")
		// Hack to get around the fact that /tmp is a symlink on
		// OS X, and it causes the directory checks to fail..
		gdir, _ = os.Getwd()
	} else {
		panic("Failed creating temporary directory")
	}
	// Make sure we get the right directory from the top level
	os.Mkdir("issues", 0755)
	b, err := New("Test Bug")
	if err != nil {
		panic("Unexpected error creating Test Bug")
	}
	t.bug = b
}
func (t *tester) Teardown() {
	os.RemoveAll(t.dir)
}
func TestTitleToDirectory(t *testing.T) {
	var assertDirectory = func(title, directory string) {
		titleStr := TitleToDir(title)
		dirStr := Directory(directory).GetShortName()

		if titleStr != dirStr {
			t.Error(fmt.Sprintf("Failed on %s: got %s but expected %s\n", title, titleStr, dirStr))
		}
	}

	assertDirectory("Test", "Test")
	assertDirectory("Test Space", "Test-Space")
	assertDirectory("Test-Dash", "Test--Dash")
	assertDirectory("Test--TripleDash", "Test---TripleDash")
	assertDirectory("Test --WithSpace", "Test_--WithSpace")
	assertDirectory("Test - What", "Test_-_What")
}

func TestNewBug(t *testing.T) {
	var gdir string
	gdir, err := ioutil.TempDir("", "newbug")
	if err == nil {
		os.Chdir(gdir)
		// Hack to get around the fact that /tmp is a symlink on
		// OS X, and it causes the directory checks to fail..
		gdir, _ = os.Getwd()
		defer os.RemoveAll(gdir)
	} else {
		t.Error("Failed creating temporary directory for detect")
		return
	}
	os.Mkdir("issues", 0755)
	b, err := New("I am a test")
	if err != nil || b == nil {
		t.Error("Unexpected error when creating New bug" + err.Error())
	}
	if b.Dir != GetIssuesDir()+TitleToDir("I am a test") {
		t.Error("Unexpected directory when creating New bug")
	}
}

func TestSetDescription(t *testing.T) {
	test := tester{}
	test.Setup()
	defer test.Teardown()

	b := test.bug

	b.SetDescription("Hello, I am a bug.")
	val, err := ioutil.ReadFile(string(b.GetDirectory()) + "/Description")
	if err != nil {
		t.Error("Could not read Description file")
	}

	if string(val) != "Hello, I am a bug." {
		t.Error("Unexpected description after SetDescription")
	}
}

func TestDescription(t *testing.T) {
	test := tester{}
	test.Setup()
	defer test.Teardown()

	b := test.bug

	desc := "I am yet another bug.\nWith Two Lines."
	b.SetDescription(desc)

	if b.Description() != desc {
		t.Error("Unexpected result from bug.Description()")
	}
}
