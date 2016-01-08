package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Create(Args ArgumentList) {
	if len(Args) < 1 || (len(Args) < 2 && Args[0] == "-n") {
		fmt.Printf("Usage: %s create [-n] Bug Description\n", os.Args[0])
		fmt.Printf("\nNo Bug Description provided.\n")
		return
	}
	var noDesc bool = false

	if Args.HasArgument("-n") {
		noDesc = true
		Args = Args[1:]
	}

	var bug bugs.Bug
	bug = bugs.Bug{
		Dir: bugs.GetIssuesDir() + bugs.TitleToDir(strings.Join(Args, " ")),
	}

	dir, _ := bug.GetDirectory()
	fmt.Printf("Created issue: %s\n", bug.Title(""))

	var mode os.FileMode
	mode = 0775
	os.Mkdir(string(dir), mode)

	if noDesc {
		txt := []byte("")
		ioutil.WriteFile(string(dir)+"/Description", txt, 0644)
	} else {
		cmd := exec.Command(getEditor(), string(dir)+"/Description")

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		if err != nil {
			log.Fatal(err)
		}
	}
}
