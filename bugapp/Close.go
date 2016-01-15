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
	for _, bugID := range args {
		if bug, err := bugs.LoadBugByHeuristic(bugID); err == nil {
			dir := bug.GetDirectory()
			fmt.Printf("Removing %s\n", dir)
			os.RemoveAll(string(dir))
		} else {
			fmt.Printf("Could not close bug %s: %s\n", bugID, err)
		}
	}
}
