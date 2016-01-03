package scm

import (
	"github.com/driusan/bug/bugs"
	"os"
	"strings"
)

type SCMNotFound struct{}

func (a SCMNotFound) Error() string {
	return "No SCM found"
}
func walkAndSearch(startpath, dirname string) string {
	if dirinfo, err := os.Stat(startpath + "/" + dirname); err == nil && dirinfo.IsDir() {
		return startpath + "/" + dirname
	}

	pieces := strings.Split(startpath, "/")

	for i := len(pieces); i > 0; i -= 1 {
		dir := strings.Join(pieces[0:i], "/")
		if dirinfo, err := os.Stat(dir + "/" + dirname); err == nil && dirinfo.IsDir() {
			return dir + "/" + dirname
		}
	}
	return ""
}

func DetectSCM() (SCMHandler, bugs.Directory, error) {
	// First look for a Git directory
	wd, _ := os.Getwd()

	gitDir := walkAndSearch(wd, ".git")
	if gitDir != "" {
		return GitManager{}, bugs.Directory(gitDir), nil
	}

	return nil, "", SCMNotFound{}
}
