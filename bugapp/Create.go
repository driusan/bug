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

	Args, argVals := Args.GetAndRemoveArguments([]string{"--tag", "--status", "--priority", "--milestone", "--identifier"})
	tag := argVals[0]
	status := argVals[1]
	priority := argVals[2]
	milestone := argVals[3]
	identifier := argVals[4]

	var bug bugs.Bug
	bug = bugs.Bug{
		Dir: bugs.GetIssuesDir() + bugs.TitleToDir(strings.Join(Args, " ")),
	}

	dir, _ := bug.GetDirectory()

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

	if tag != "" {
		bug.TagBug(bugs.Tag(tag))
	}
	if status != "" {
		bug.SetStatus(status)
	}
	if priority != "" {
		bug.SetPriority(priority)
	}
	if milestone != "" {
		bug.SetMilestone(milestone)
	}
	if identifier != "" {
		bug.SetIdentifier(identifier)
	}
	fmt.Printf("Created issue: %s\n", bug.Title(""))
}
