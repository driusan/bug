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

func walkAndSearch(startpath string, dirnames []string) (fullpath, scmtype string) {
	for _, scmtype := range dirnames {
		if dirinfo, err := os.Stat(startpath + "/" + scmtype); err == nil && dirinfo.IsDir() {
			return startpath + "/" + scmtype, scmtype
		}
	}

	pieces := strings.Split(startpath, "/")

	for i := len(pieces); i > 0; i -= 1 {
		dir := strings.Join(pieces[0:i], "/")
		for _, scmtype := range dirnames {
			if dirinfo, err := os.Stat(dir + "/" + scmtype); err == nil && dirinfo.IsDir() {
				return dir + "/" + scmtype, scmtype
			}
		}
	}
	return "", ""
}

func DetectSCM(options map[string]bool) (SCMHandler, bugs.Directory, error) {
	// First look for a Git directory
	wd, _ := os.Getwd()

	dirFound, scmtype := walkAndSearch(wd, []string{".git", ".hg"})
	if dirFound != "" && scmtype == ".git" {
		var gm GitManager
		if val, ok := options["autoclose"]; ok {
			gm.Autoclose = val
		}
		if val, ok := options["use_bug_prefix"]; ok {
			gm.UseBugPrefix = val
		}
		return gm, bugs.Directory(dirFound), nil
	}
	if dirFound != "" && scmtype == ".hg" {
		return HgManager{}, bugs.Directory(dirFound), nil
	}

	return nil, "", SCMNotFound{}
}
