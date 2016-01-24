package bugs

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestDirectoryToTitle(t *testing.T) {
	var assertTitle = func(directory, title string) {
		dir := Directory(directory)
		if dir.ToTitle() != title {
			t.Error("Failed on " + directory + ": got " + dir.ToTitle() + " but expected " + title)
		}
	}
	assertTitle("Test", "Test")
	assertTitle("Test-Multiword", "Test Multiword")
	assertTitle("Test--Dash", "Test-Dash")
	assertTitle("Test---Dash", "Test--Dash")
	assertTitle("Test_--TripleDash", "Test --TripleDash")
	assertTitle("Test_-_What", "Test - What")
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
