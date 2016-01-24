package bugs

import (
    "time"
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
	root := GetRootDir()
	if root == "" {
		return root
	}
	return GetRootDir() + "/issues/"
}

type Directory string

func (d Directory) GetShortName() Directory {
	pieces := strings.Split(string(d), "/")
	return Directory(pieces[len(pieces)-1])
}

func (d Directory) ToTitle() string {
	multidash := regexp.MustCompile("([_]*)-([-_]*)")
	dashReplacement := strings.Replace(string(d), " ", "/", -1)
	return multidash.ReplaceAllStringFunc(dashReplacement, func(match string) string {
		if match == "-" {
			return " "
		}
		if strings.Count(match, "_") == 0 {
			return match[1:]
		}
		return strings.Replace(match, "_", " ", -1)
	})
}

func (d Directory) LastModified() time.Time {
    var t time.Time
    stat, err := os.Stat(string(d))
    if err != nil {
        panic("Directory " + string(d) + " is not a directory.")
    }

    if stat.IsDir() == false {
        return stat.ModTime()
    }

    dir, _ := os.Open(string(d))
    files, _ := dir.Readdir(-1)
    if len(files) == 0 {
        t = stat.ModTime()
    }
    for _, file := range files {
        if file.IsDir() {
            mtime := (d +"/" + Directory(file.Name())).LastModified()
            if mtime.After(t) {
                t = mtime
            }
        } else {
            mtime := file.ModTime()
            if mtime.After(t) {
                t = mtime
            }

        }
    }
    return t
}
