package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Edit(args ArgumentList) {
	issues, _ := ioutil.ReadDir(string(bugs.GetIssuesDir()))

	var file string
	var idx int
	var err error
	switch len(args) {
	case 1:
		// If there's only 1 argument, it's an issue
		// identifier and it's editing the description.
		// So set the variables and fallthrough to the
		// 2 argument (editing a specific fieldname)
		// case
		idx, err = strconv.Atoi(args[0])
		file = "Description"
		fallthrough
	case 2:
		// If there's exactly 2 arguments, idx and
		// file haven't been set by the first case
		// statement, so set them, but everything else
		// is the same
		if len(args) == 2 {
			idx, err = strconv.Atoi(args[1])
			file = args[0]
		}

		if idx > len(issues) || idx < 1 {
			fmt.Printf("Invalid issue number %d\n", idx)
			return
		}
		dir := bugs.Directory(bugs.GetIssuesDir()) + bugs.Directory(issues[idx-1].Name())

		switch title := strings.Title(file); title {
		case "Milestone", "Status", "Priority", "Identifier":
			file = title
		}
		cmd := exec.Command(getEditor(), string(dir)+"/"+file)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Printf("Usage: %s edit [fieldname] issuenum\n", os.Args[0])
		fmt.Printf("\nNo issue number specified\n")
	}
}
