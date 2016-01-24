package bugs

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGetRootDirWithEnvironmentVariable(t *testing.T) {
	var gdir string
	gdir, err := ioutil.TempDir("", "rootdirbug")
	if err == nil {
		os.Chdir(gdir)
		// Hack to get around the fact that /tmp is a symlink on
		// OS X, and it causes the directory checks to fail..
		gdir, _ = os.Getwd()
		defer os.RemoveAll(gdir)
	} else {
		t.Error("Failed creating temporary directory")
		return
	}
	os.Mkdir("issues", 0755)
	os.Setenv("PMIT", "/tmp/abc")
	dir := GetRootDir()
	if dir != Directory("/tmp/abc") {
		t.Error("Did not get proper directory according to environment variable")
	}
}
func TestGetRootDirFromDirectoryTree(t *testing.T) {
	var gdir string
	gdir, err := ioutil.TempDir("", "rootdirbug")
	if err == nil {
		os.Chdir(gdir)
		os.Unsetenv("PMIT")
		// Hack to get around the fact that /tmp is a symlink on
		// OS X, and it causes the directory checks to fail..
		gdir, _ = os.Getwd()
		defer os.RemoveAll(gdir)
	} else {
		t.Error("Failed creating temporary directory")
		return
	}
	// Make sure we get the right directory from the top level
	os.Mkdir("issues", 0755)
	dir := GetRootDir()
	if dir != Directory(gdir) {
		t.Error("Did not get proper directory according to walking the tree:" + dir)
	}
	// Now go deeper into the tree and try the same thing..
	os.Mkdir("abc/123", 0755)
	os.Chdir("abc/123")
	dir = GetRootDir()
	if dir != Directory(gdir) {
		t.Error("Did not get proper directory according to walking the tree:" + dir)
	}
}

func TestGetIssuesDir(t *testing.T) {
	os.Setenv("PMIT", "/tmp/abc")
	dir := GetIssuesDir()
	if dir != "/tmp/abc/issues/" {
		t.Error("Did not get correct issues directory")
	}
}
func TestShortName(t *testing.T) {
	var dir Directory = "/hello/i/am/a/test"
	if short := dir.GetShortName(); short != Directory("test") {
		t.Error("Unexpected short name: " + string(short))
	}
}
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
