package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"github.com/driusan/bug/scm"
)

func Commit(args ArgumentList) {
	options := make(map[string]bool)
	if !args.HasArgument("--no-autoclose") {
		options["autoclose"] = true
	} else {
		options["autoclose"] = false
	}

	scm, _, err := scm.DetectSCM(options)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	err = scm.Commit(bugs.GetIssuesDir(), "Added or removed issues with the tool \"bug\"")

	if err != nil {
		fmt.Printf("Could not commit: %s\n", err.Error())
		return
	}
}
