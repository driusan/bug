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
		fmt.Printf("Usage: " + os.Args[0] + " create [-n] Issue Title\n\n")
		fmt.Printf("This will create an issue with the title Issue Title\n\n")
		fmt.Printf("An editor will be opened automatically for you to enter\n")
		fmt.Printf("a more detailed description.\n\n")
		fmt.Printf("If your EDITOR environment variable is set, it will be\n")
		fmt.Printf("used, otherwise the default is vim.\n")
		fmt.Printf("If the first argument to create is \"-n\", then " + os.Args[0] + " will not open any editor and create an empty Description\n\n")
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
	case "commit":
		fmt.Printf("Usage: " + os.Args[0] + " commit\n\n")
		fmt.Printf(`This will commit any new, modified, or removed issues to
git.

Your working tree and staging area should be otherwise unaffected by using
this command.
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
    case "tag":
		fmt.Printf("Usage: " + os.Args[0] + " tag IssueNumber [tags]\n\n")
    fmt.Printf(`This will tag the given IssueNumber with the tags
given as parameters. At least one tag is required.

Tags can be any string which would make a valid file name.
`);
	case "help":
		fallthrough
	default:
		fmt.Printf("Usage: " + os.Args[0] + " command [options]\n\n")
		fmt.Printf("Use \"bug help [command]\" for more information about any command below\n\n")
		fmt.Printf("Valid commands\n")
		fmt.Printf("\tcreate\tFile a new bug\n")
		fmt.Printf("\tlist\tList existing bugs\n")
		fmt.Printf("\tedit\tEdit an existing bug\n")
		fmt.Printf("\ttag\tTag a bug with a category\n")
		fmt.Printf("\tclose\tDelete an existing bug\n")
		fmt.Printf("\tcommit\tCommit any new, changed or deleted bug to git\n")
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

func (a BugApplication) Tag(Args []string) {
	if len(Args) < 2 {
		fmt.Printf("Invalid usage. Must provide issue and tags\n.")
		return
	}

	issues, err := ioutil.ReadDir(getRootDir() + "/issues")
	if err != nil {
		fmt.Printf("Unknown error reading directory: %s\n", err.Error())
		return
	}
	idx, err := strconv.Atoi(Args[0])
    idx = idx - 1
	if err != nil {
		fmt.Printf("Unknown looking up bug: %s\n", err)
		return
	}
    if idx >= len(issues) || idx < 0 {
		fmt.Printf("Invalid issue index.\n")
        return
    }
	var b Bug
	b.LoadBug(Directory(getRootDir() + "/issues/" + issues[idx].Name()))
	for _, tag := range Args[1:] {
        b.TagBug(tag)
	}

}
func (a BugApplication) Create(Args []string) {
	var bug Bug
	var noDesc bool = false

	if Args != nil && Args[0] == "-n" {
		noDesc = true
		Args = Args[1:]
	}
	bug = Bug{
		Title: strings.Join(Args, " "),
	}

	dir, _ := bug.GetDirectory()
	fmt.Printf("Created issue: %s\n", bug.Title)

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

func (a BugApplication) Dir() {
	fmt.Printf("%s", getRootDir()+"/issues")
}

// This will try and commit the $(bug pwd) directory
// transparently. It does the following steps:
//
// 1. "git stash create"
// 2. "git reset --mixed" (unstage the user's currently staged files)
// 3. "git add $(bug pwd)"
// 4. "git commit"
// 5a. "git reset --hard" (if there was any stash created,
// 						this is necessary for 5b to work.)
// 5b. "git stash apply --index" the stash from step 1
func (a BugApplication) Commit() {
	cmd := exec.Command("git", "stash", "create")

	output, err := cmd.Output()

	if err != nil {
		fmt.Printf("Could not execute git stash create")
	}
	var stashHash string = strings.Trim(string(output), "\n")

	// Unstage everything, if there was anything stashed, so that
	// we don't commit things that the user has staged that aren't
	// issues
	if stashHash != "" {
		cmd = exec.Command("git", "reset", "--mixed")
		err = cmd.Run()

		if err != nil {
		}
	}

	// Commit the issues directory
	// git add $(bug pwd)
	// git commit -m "Added new issues" -q
	cmd = exec.Command("git", "add", "-A", getRootDir()+"/issues")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Could not add to index?\n")
	}
	cmd = exec.Command("git", "commit", "-m", "Added or removes issues with the tool \"bug\"", "-q")
	err = cmd.Run()
	if err != nil {
		// If nothing was added commit will have an error,
		// but we don't care it just means there's nothing
		// to commit.
		fmt.Printf("No new issues commited\n")
	}

	// There were changes that had been stashed, so we need
	// to restore them with git stash apply.. first, we
	// need to do a "git reset --hard" so that the dirty working
	// tree doesn't cause an error. This isn't as scary as it
	// sounds, since immediately after git reset --hard we apply
	// a stash which has the exact same changes that we just threw
	// away.
	if stashHash != "" {
		cmd = exec.Command("git", "reset", "--hard")
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error resetting the git working tree")
			// If nothing was stashed, it's not the end of the world.
			//fmt.Printf("Could not pop from stash\n")
		}
		cmd = exec.Command("git", "stash", "apply", "--index", stashHash)
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error resetting the git working tree")
			// If nothing was stashed, it's not the end of the world.
			//fmt.Printf("Could not pop from stash\n")
		}
	}
}
