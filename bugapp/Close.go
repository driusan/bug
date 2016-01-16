package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"os"
)

func Close(args ArgumentList) {
	// No parameters, print a list of all bugs
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s close BugID\n\nMust provide an BugID to close as parameter\n", os.Args[0])
		return
	}

	// There were parameters, so show the full description of each
	// of those issues
	var bugsToClose []string
	for _, bugID := range args {
		if bug, err := bugs.LoadBugByHeuristic(bugID); err == nil {
			dir := bug.GetDirectory()
			bugsToClose = append(bugsToClose, string(dir))
		} else {
			fmt.Fprintf(os.Stderr, "Could not close bug %s: %s\n", bugID, err)
		}
	}
	for _, dir := range bugsToClose {
		fmt.Printf("Removing %s\n", dir)
		os.RemoveAll(dir)
	}
}
