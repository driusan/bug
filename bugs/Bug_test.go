package bugs

import (
	"fmt"
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
