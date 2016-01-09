package bugapp

import (
	"crypto/sha1"
	"fmt"
	"github.com/driusan/bug/bugs"
	"os"
	"strconv"
	"strings"
)

func Identifier(args ArgumentList) {
	if len(args) < 1 {
		fmt.Printf("Usage: %s identifier issuenum [value]\n", os.Args[0])
		return
	}

	idx, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Invalid issue number. \"%s\" is not a number.\n\n", args[0])
		fmt.Printf("Usage: %s identifier issuenum [value]\n", os.Args[0])
		return
	}

	b, err := bugs.LoadBugByIndex(idx)
	if err != nil {
		fmt.Printf("Invalid issue number %s\n", args[0])
		return
	}
	if len(args) > 1 {
		var newValue string
		if args.HasArgument("--generate") {
			hash := sha1.Sum([]byte(b.Title("")))
			newValue = fmt.Sprintf("b%x", hash)[0:5]
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
