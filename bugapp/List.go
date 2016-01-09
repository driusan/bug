package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"os"
	"strconv"
)

func listTags(files []os.FileInfo, args ArgumentList) {
	b := bugs.Bug{}
	for idx, _ := range files {
		b.LoadBug(bugs.Directory(bugs.GetIssuesDir() + bugs.Directory(files[idx].Name())))

		for _, tag := range args {
			if b.HasTag(bugs.Tag(tag)) {
				fmt.Printf("Issue %d: %s\n", idx+1, b.Title("tags"))
			}
		}
	}
}
func List(args ArgumentList) {
	issues, _ := ioutil.ReadDir(string(bugs.GetIssuesDir()))

	var wantTags bool = false
	if args.HasArgument("--tags") {
		wantTags = true
	}

	// No parameters, print a list of all bugs
	if len(args) == 0 || (wantTags && len(args) == 1) {
		for idx, issue := range issues {
			if issue.IsDir() != true {
				continue
			}
			var dir bugs.Directory = bugs.GetIssuesDir() + bugs.Directory(issue.Name())
			b := bugs.Bug{dir}
			var name string
			if id := b.Identifier(); id != "" {
				name = fmt.Sprintf("Issue %s", id)
			} else {
				name = fmt.Sprintf("Issue %d", idx+1)
			}
			if wantTags == false {
				fmt.Printf("%s: %s\n", name, b.Title(""))
			} else {
				fmt.Printf("%s: %s\n", name, b.Title("tags"))
			}
		}
		return
	}

	// There were parameters, so show the full description of each
	// of those issues
	b := bugs.Bug{}
	for i, length := 0, len(args); i < length; i += 1 {
		idx, err := strconv.Atoi(args[i])
		if err != nil {
			listTags(issues, args)
			return
		}
		if idx > len(issues) || idx < 1 {
			fmt.Printf("Invalid issue number %d\n", idx)
			continue
		}
		if err == nil {
			b.LoadBug(bugs.Directory(bugs.GetIssuesDir() + bugs.Directory(issues[idx-1].Name())))
			b.ViewBug()
			if i < length-1 {
				fmt.Printf("\n--\n\n")
			}
		}
	}
	fmt.Printf("\n")
}
