package bugs

import (
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
}

func TestBugGetDirectory(t *testing.T) {
	b := Bug{}
	var assertDirectory = func(title, directory string) {
		b.Title = title
		dir, _ := b.GetDirectory()
		dirStr := string(dir.GetShortName())

		if directory != dirStr {
			t.Error("Failed on " + title + ": got " + dirStr + " but expected " + directory)
		}
	}
	assertDirectory("Test", "Test")
	assertDirectory("Test Space", "Test-Space")
	assertDirectory("Test-Dash", "Test--Dash")
	assertDirectory("Test--TripleDash", "Test---TripleDash")

}
