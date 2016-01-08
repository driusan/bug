package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"os"
	"strconv"
)

func Close(args ArgumentList) {
	issues, _ := ioutil.ReadDir(string(bugs.GetIssuesDir()))

	// No parameters, print a list of all bugs
	if len(args) == 0 {
		fmt.Printf("Usage: %s close IssueNumber\n\nMust provide an IssueNumber to close as parameter\n", os.Args[0])
		return
	}

	// There were parameters, so show the full description of each
	// of those issues
	for i := 0; i < len(args); i += 1 {
		idx, err := strconv.Atoi(args[i])
		if idx > len(issues) || idx < 1 {
			fmt.Printf("Invalid issue number %d\n", idx)
			continue
		}
		if err == nil {
			dir := bugs.GetIssuesDir() + bugs.Directory(issues[idx-1].Name())
			fmt.Printf("Removing %s\n", dir)
			os.RemoveAll(string(dir))
		}
	}
}
