package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"github.com/driusan/bug/scm"
)

func Purge() {
	scm, _, err := scm.DetectSCM(make(map[string]bool))

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	err = scm.Purge(bugs.GetIssuesDir())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
}
