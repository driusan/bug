package scm

import (
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"os"
	"testing"
)

func TestDetectGit(t *testing.T) {
	var gdir string
	gdir, err := ioutil.TempDir("", "gitdetect")
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
	// Fake a git repo
	os.Mkdir(".git", 0755)

	options := make(map[string]bool)
	handler, dir, err := DetectSCM(options)
	if err != nil {
		t.Error("Unexpected while detecting repo type: " + err.Error())
	}
	if dir != bugs.Directory(gdir+"/.git") {
		t.Error("Unexpected directory found when trying to detect git repo" + dir)
	}
	switch handler.(type) {
	case GitManager:
		// GitManager is what we expect, don't fall through
		// to the error
	default:
		t.Error("Unexpected SCMHandler found for Git")
	}

	// Go somewhere higher in the tree and do it again
	os.MkdirAll("tmp/abc/hello", 0755)
	os.Chdir("tmp/abc/hello")
	handler, dir, err = DetectSCM(options)
	if err != nil {
		t.Error("Unexpected while detecting repo type: " + err.Error())
	}
	if dir != bugs.Directory(gdir+"/.git") {
		t.Error("Unexpected directory found when trying to detect git repo" + dir)
	}
	switch handler.(type) {
	case GitManager:
		// GitManager is what we expect, don't fall through
		// to the error
	default:
		t.Error("Unexpected SCMHandler found for Git")
	}
}

func TestDetectHg(t *testing.T) {
	var gdir string
	gdir, err := ioutil.TempDir("", "hgdetect")
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
	// Fake a git repo
	os.Mkdir(".hg", 0755)

	options := make(map[string]bool)
	handler, dir, err := DetectSCM(options)
	if err != nil {
		t.Error("Unexpected while detecting repo type: " + err.Error())
	}
	if dir != bugs.Directory(gdir+"/.hg") {
		t.Error("Unexpected directory found when trying to detect git repo" + dir)
	}
	switch handler.(type) {
	case HgManager:
		// HgManager is what we expect, don't fall through
		// to the error
	default:
		t.Error("Unexpected SCMHandler found for Mercurial")
	}

	// Go somewhere higher in the tree and do it again
	os.MkdirAll("tmp/abc/hello", 0755)
	os.Chdir("tmp/abc/hello")
	handler, dir, err = DetectSCM(options)
	if err != nil {
		t.Error("Unexpected while detecting repo type: " + err.Error())
	}
	if dir != bugs.Directory(gdir+"/.hg") {
		t.Error("Unexpected directory found when trying to detect git repo" + dir)
	}
	switch handler.(type) {
	case HgManager:
		// HgManager is what we expect, don't fall through
		// to the error
	default:
		t.Error("Unexpected SCMHandler found for Mercurial")
	}
}
