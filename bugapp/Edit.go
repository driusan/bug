package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Edit(args ArgumentList) {

	var file, bugID string
	switch len(args) {
	case 1:
		// If there's only 1 argument, it's an issue
		// identifier and it's editing the description.
		// So set the variables and fallthrough to the
		// 2 argument (editing a specific fieldname)
		// case
		bugID = args[0]
		file = "Description"
		fallthrough
	case 2:
		// If there's exactly 2 arguments, idx and
		// file haven't been set by the first case
		// statement, so set them, but everything else
		// is the same
		if len(args) == 2 {
			bugID = args[1]
			file = args[0]
		}

		b, err := bugs.LoadBugByHeuristic(bugID)
		if err != nil {
			fmt.Printf("Invalid BugID %s\n", bugID)
			return
		}

		dir := b.GetDirectory()

		switch title := strings.Title(file); title {
		case "Milestone", "Status", "Priority", "Identifier":
			file = title
		}
		fmt.Printf("Launching in %s/%s", dir, file)
		cmd := exec.Command(getEditor(), string(dir)+"/"+file)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Printf("Usage: %s edit [fieldname] BugID\n", os.Args[0])
		fmt.Printf("\nNo BugID specified\n")
	}
}
