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

func (a BugApplication) Help(args ...string) {
	var cmd string
	if args == nil {
		cmd = "help"

	}
	if len(args) <= 1 {
		cmd = "help"
	} else {
		cmd = args[1]
	}
	switch cmd {
	case "create":
		fmt.Printf("Usage: " + os.Args[0] + " create Issue Title\n\n")
		fmt.Printf("This will create an issue with the title Issue Title\n\n")
		fmt.Printf("An editor will be opened automatically for you to enter\n")
		fmt.Printf("a more detailed description.\n\n")
		fmt.Printf("If your EDITOR environment variable is set, it will be\n")
		fmt.Printf("used, otherwise the default is vim.\n")
	case "list":
		fmt.Printf("Usage: " + os.Args[0] + " list [issue numbers]\n\n")
		fmt.Printf(`This will list the issues found in the current environment

With no arguments, titles will be printed to the screen along with the issue
number that can be used to reference this issue on the command line.

If 1 or more issue numbers are provided, the whole issue including description
will be printed to stdout.
`)

	case "edit":
		fmt.Printf("Usage: " + os.Args[0] + " edit IssueNumber\n\n")
		fmt.Printf(
			`This will launch your standard editor to edit the description of the bug numbered 
IssueNumber, where IssueNumber is a reference to same index provided with a
"bug list" command.
`)
	case "rm":
		fallthrough
	case "close":
		fmt.Printf("Usage: " + os.Args[0] + " close IssueNumber\n\n")
		fmt.Printf(`This will delete the issue numbered IssueNumber. IssueNumbers
correspond to the number in the "bug list" command.

Note that closing a bug will cause all existing bugs to be be renumbered and
IssueNumbers are not intended to be stable.

Also note that this does not remove the issue from git, but only from the file
system. If you want to remove an issue that is tracked by git, you'll have to
manually "git rm -r" the directory from the directory that's printed to the
screen if you execute "bug dir".`)
		fmt.Printf("\n\n\"%s rm\" is an alias for this \"%s close\"\n", os.Args[0], os.Args[0])
	case "purge":
		fmt.Printf("Usage: " + os.Args[0] + " purge\n\n")
		fmt.Printf(`This will delete any bugs that are not currently tracked by
git.

It is an alias for "git clean -fd $(bug dir)"
`)
	case "env":
		fmt.Printf("Usage: " + os.Args[0] + " env\n\n")
		fmt.Printf(`This will print the environment used by the bug command to stdout.

Use this command if you want to see what directory bug create is
using to store bugs, or what editor will be invoked by bug create/edit.
`)

	case "pwd":
		fallthrough
	case "dir":
		fmt.Printf("Usage: " + os.Args[0] + " dir\n\n")
		fmt.Printf(`This will bug directory to stdout, so you can use it as a subcommand
for arguments to any arbitrary shell commands. For example "cd $(bug dir)" or 
"git rm -r $(bug dir)/Issue-Title"
`)
		fmt.Printf("\n\n\"%s pwd\" is an alias for this \"%s dir\"\n", os.Args[0], os.Args[0])
	case "help":
		fallthrough
	default:
		fmt.Printf("Usage: " + os.Args[0] + " command [options]\n\n")
		fmt.Printf("Use \"bug help [command]\" for more information about any command below\n\n")
		fmt.Printf("Valid commands\n")
		fmt.Printf("\tcreate\tFile a new bug\n")
		fmt.Printf("\tlist\tList existing bugs\n")
		fmt.Printf("\tedit\tEdit an existing bug\n")
		fmt.Printf("\tclose\tDelete an existing bug\n")
		fmt.Printf("\tpurge\tRemove all issues not tracked by git\n")
		fmt.Printf("\trm\tAlias of close\n")
		fmt.Printf("\tenv\tShow settings that bug will use if invoked from this directory\n")
		fmt.Printf("\tdir\tPrints the issues directory to stdout (useful subcommand in the shell)\n")
		fmt.Printf("\tpwd\tAlias of dir\n")
		fmt.Printf("\thelp\tShow this screen\n")
	}
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

func (a BugApplication) Edit(args []string) {
	issues, _ := ioutil.ReadDir(getRootDir() + "/issues")

	// No parameters, print a list of all bugs
	if len(args) == 1 {
		idx, err := strconv.Atoi(args[0])
		if idx > len(issues) || idx < 1 {
			fmt.Printf("Invalid issue number %d\n", idx)
			return
		}
		dir := Directory(getRootDir() + "/issues/" + issues[idx-1].Name())
		cmd := exec.Command(getEditor(), string(dir)+"/Description")

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("Usage: %s edit issuenum\n", os.Args[0])
		fmt.Printf("\nNo issue number specified\n")
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
	fmt.Printf("Created issue: %s\n", bug.Title)

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

func (a BugApplication) Commit() {
	cmd := exec.Command("git", "stash", "save")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	cmd = exec.Command("git", "add", getRootDir()+"/issues")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	cmd = exec.Command("git", "commit", "-m", "Added new issues", "-q")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("No issues commited\n")
		//log.Fatal(err)
	}
	cmd = exec.Command("git", "stash", "pop")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
