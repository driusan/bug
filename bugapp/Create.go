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
		fmt.Fprintf(os.Stderr, "Usage: %s create [-n] Bug Description\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nNo Bug Description provided.\n")
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

	if Args.HasArgument("--generate-id") {
		for i, token := range Args {
			if token == "--generate-id" {
				if i+1 < len(Args) {
					Args = append(Args[:i], Args[i+1:]...)
					break
				} else {
					Args = Args[:i]
					break
				}
			}
		}
		identifier = generateID(strings.Join(Args, " "))
	}

	// It's possible there were arguments provided, but still no title
	// included. Do another check before trying to create the bug.
	if strings.TrimSpace(strings.Join(Args, " ")) == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s create [-n] Bug Description\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nNo Bug Description provided.\n")
		return
	}
	var bug bugs.Bug
	bug = bugs.Bug{
		Dir: bugs.GetIssuesDir() + bugs.TitleToDir(strings.Join(Args, " ")),
	}

	dir := bug.GetDirectory()

	var mode os.FileMode
	mode = 0775
	err := os.Mkdir(string(dir), mode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%s error: mkdir\n", os.Args[0])
		log.Fatal(err)
	}

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
