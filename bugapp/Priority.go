package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"os"
	"strings"
)

func fieldHandler(command string, args ArgumentList,
	setCallback func(bugs.Bug, string) error, retrieveCallback func(bugs.Bug) string) {
	if len(args) < 1 {
		fmt.Printf("Usage: %s %s BugID [set %s]\n", os.Args[0], command, command)
		return
	}

	b, err := bugs.LoadBugByHeuristic(args[0])
	if err != nil {
		fmt.Printf("Invalid BugID: %s\n", err.Error())
		return
	}
	if len(args) > 1 {
		newValue := strings.Join(args[1:], " ")
		err := setCallback(*b, newValue)
		if err != nil {
			fmt.Printf("Error setting %s: %s", command, err.Error())
		}
	} else {
		val := retrieveCallback(*b)
		if val == "" {
			fmt.Printf("%s not defined\n", command)
		} else {
			fmt.Printf("%s\n", val)
		}
	}
}
func Priority(args ArgumentList) {
	fieldHandler("priority", args, bugs.Bug.SetPriority, bugs.Bug.Priority)
}
