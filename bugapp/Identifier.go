package bugapp

import (
	"crypto/sha1"
	"fmt"
	"github.com/driusan/bug/bugs"
	"os"
	"strings"
)

func generateID(val string) string {
	hash := sha1.Sum([]byte(val))
	return fmt.Sprintf("b%x", hash)[0:5]
}
func Identifier(args ArgumentList) {
	if len(args) < 1 {
		fmt.Printf("Usage: %s identifier BugID [value]\n", os.Args[0])
		return
	}

	b, err := bugs.LoadBugByHeuristic(args[0])
	if err != nil {
		fmt.Printf("Invalid BugID: %s\n", err.Error())
		return
	}
	if len(args) > 1 {
		var newValue string
		if args.HasArgument("--generate") {
			newValue = generateID(b.Title(""))
			fmt.Printf("Generated id %s for bug\n", newValue)
		} else {
			newValue = strings.Join(args[1:], " ")
		}
		err := b.SetIdentifier(newValue)
		if err != nil {
			fmt.Printf("Error setting identifier: %s", err.Error())
		}
	} else {
		val := b.Identifier()
		if val == "" {
			fmt.Printf("Identifier not defined\n")
		} else {
			fmt.Printf("%s\n", val)
		}
	}
}
