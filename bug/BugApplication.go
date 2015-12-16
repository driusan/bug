package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	//"regex"
	"strings"
)

type BugApplication struct{}

func (a BugApplication) Help() {
	fmt.Printf("Usage: " + os.Args[0] + " command [options]\n\n")
	fmt.Printf("Valid commands\n")
	fmt.Printf("\tcreate\tFile a new bug\n")
	fmt.Printf("\tlist\tList existing bugs\n")
	fmt.Printf("\tclose\tDelete an existing bug\n")
	fmt.Printf("\tpurge\tRemove all issues not tracked by git\n")
	fmt.Printf("\trm\tAlias of close\n")
	fmt.Printf("\tenv\tShow settings that bug will use if invoked from this directory\n")
	fmt.Printf("\tdir\tPrints the issues directory to stdout (useful subcommand in the shell)\n")
	fmt.Printf("\thelp\tShow this screen\n")
}

func (a BugApplication) Env() {
	fmt.Printf("Settings used by this command:\n")
	fmt.Printf("\nIssues directory:\t%s/issues/", getRootDir())
	fmt.Printf("\nEditor:\t%s", getEditor())
	fmt.Printf("\n")
}

func (a BugApplication) List(args []string) {
	issues, _ := ioutil.ReadDir(getRootDir() + "/issues")

	// No parameters, print a list of all bugs
	if len(args) == 0 {
		for idx, issue := range issues {
			var dir Directory = Directory(issue.Name())
			fmt.Printf("Issue %d: %s\n", idx+1, dir.ToTitle())
		}
		return
	}

	// There were parameters, so show the full description of each
	// of those issues
	b := Bug{}
	for i := 0; i < len(args); i += 1 {
		idx, err := strconv.Atoi(args[i])
		if idx > len(issues) || idx < 1 {
			fmt.Printf("Invalid issue number %d\n", idx)
			continue
		}
		if err == nil {
			b.LoadBug(Directory(getRootDir() + "/issues/" + issues[idx-1].Name()))
			b.ViewBug()
		}
	}
}
func (a BugApplication) Close(args []string) {
	issues, _ := ioutil.ReadDir(getRootDir() + "/issues")

	// No parameters, print a list of all bugs
	if len(args) == 0 {
		fmt.Printf("Must provide bug to close as parameter\n")
		return
	}

	// There were parameters, so show the full description of each
	// of those issues
	for i := 0; i < len(args); i += 1 {
		idx, err := strconv.Atoi(args[i])
		if idx > len(issues) || idx < 1 {
			fmt.Printf("Invalid issue number %d\n", idx)
			continue
		}
		if err == nil {
			var dir string = getRootDir() + "/issues/" + issues[idx-1].Name()
			fmt.Printf("Removing %s\n", dir)
			os.RemoveAll(dir)
		}
	}
}

func (a BugApplication) Purge() {
	cmd := exec.Command("git", "clean", "-fd", getRootDir()+"/issues")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
}
func (a BugApplication) Create(Args []string) {
	var bug Bug

	bug = Bug{
		Title: strings.Join(Args, " "),
	}

	dir, _ := bug.GetDirectory()
	fmt.Printf("Created a bug: %s\n\tIn directory: %s", bug.Title, dir)

	var mode os.FileMode
	mode = 0775
	os.Mkdir(string(dir), mode)

	cmd := exec.Command(getEditor(), string(dir)+"/Description")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func (a BugApplication) Dir() {
	fmt.Printf("%s", getRootDir()+"/issues")
}
