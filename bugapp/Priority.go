package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"os"
	"strconv"
	"strings"
)

func fieldHandler(command string, args ArgumentList,
	setCallback func(bugs.Bug, string) error, retrieveCallback func(bugs.Bug) string) {
	if len(args) < 1 {
		fmt.Printf("Usage: %s %s issuenum [set %s]\n", os.Args[0], command, command)
		return
	}

	idx, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Invalid issue number. \"%s\" is not a number.\n\n", args[0])
		fmt.Printf("Usage: %s %s issuenum [set %s]\n", os.Args[0], command, command)
		return
	}

	b, err := bugs.LoadBugByIndex(idx)
	if err != nil {
		fmt.Printf("Invalid issue number %s\n", args[0])
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
