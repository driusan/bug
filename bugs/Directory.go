package bugs

import (
	"os"
	"strings"
)

func GetRootDir() string {
	dir := os.Getenv("PMIT")
	if dir != "" {
		return dir
	}

	wd, _ := os.Getwd()

	if dirinfo, err := os.Stat(wd + "/issues"); err == nil && dirinfo.IsDir() {
		return wd
	}

	// There's no environment variable and no issues
	// directory, so walk up the tree until we find one
	pieces := strings.Split(wd, "/")

	for i := len(pieces); i > 0; i -= 1 {
		dir := strings.Join(pieces[0:i], "/")
		if dirinfo, err := os.Stat(dir + "/issues"); err == nil && dirinfo.IsDir() {
			return dir
		}
	}
	return ""
}

type Directory string

func (d Directory) GetShortName() Directory {
	pieces := strings.Split(string(d), "/")
	return Directory(pieces[len(pieces)-1])
}

func (d Directory) ToTitle() string {
	tokens := strings.Split(string(d), "-")
	return strings.Join(tokens, " ")
}
