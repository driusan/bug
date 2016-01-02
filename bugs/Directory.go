package bugs

import (
	"os"
	"regexp"
	"strings"
)

func GetRootDir() Directory {
	dir := os.Getenv("PMIT")
	if dir != "" {
		return Directory(dir)
	}

	wd, _ := os.Getwd()

	if dirinfo, err := os.Stat(wd + "/issues"); err == nil && dirinfo.IsDir() {
		return Directory(wd)
	}

	// There's no environment variable and no issues
	// directory, so walk up the tree until we find one
	pieces := strings.Split(wd, "/")

	for i := len(pieces); i > 0; i -= 1 {
		dir := strings.Join(pieces[0:i], "/")
		if dirinfo, err := os.Stat(dir + "/issues"); err == nil && dirinfo.IsDir() {
			return Directory(dir)
		}
	}
	return ""
}

func GetIssuesDir() Directory {
	return GetRootDir() + "/issues/"
}

type Directory string

func (d Directory) GetShortName() Directory {
	pieces := strings.Split(string(d), "/")
	return Directory(pieces[len(pieces)-1])
}

func (d Directory) ToTitle() string {
	multidash := regexp.MustCompile("-(-*)")
	return multidash.ReplaceAllStringFunc(string(d), func(match string) string {
		if match == "-" {
			return " "
		}
		return match[1:]
	})
}
