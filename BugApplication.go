package main

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

type ArgumentList []string

func (args ArgumentList) HasArgument(arg string) bool {
	for _, argCandidate := range args {
		if arg == argCandidate {
			return true
		}
	}
	return false
}

type BugApplication struct{}

func (a BugApplication) Env() {
	fmt.Printf("Settings used by this command:\n")
	fmt.Printf("\nIssues directory:\t%s\n", bugs.GetIssuesDir())
	fmt.Printf("\nEditor:\t%s", getEditor())
	fmt.Printf("\n")
}

func listTags(files []os.FileInfo, args ArgumentList) {
	b := bugs.Bug{}
	for idx, _ := range files {
		b.LoadBug(bugs.Directory(bugs.GetIssuesDir() + bugs.Directory(files[idx].Name())))

		for _, tag := range args {
			if b.HasTag(bugs.Tag(tag)) {
				fmt.Printf("Issue %d: %s\n", idx+1, b.Title("tags"))
			}
		}
	}
}
func (a BugApplication) List(args ArgumentList) {
	issues, _ := ioutil.ReadDir(string(bugs.GetIssuesDir()))

	// No parameters, print a list of all bugs
	if len(args) == 0 {
		for idx, issue := range issues {
			var dir bugs.Directory = bugs.Directory(issue.Name())
			fmt.Printf("Issue %d: %s\n", idx+1, dir.ToTitle())
		}
		return
	}

	// There were parameters, so show the full description of each
	// of those issues
	b := bugs.Bug{}
	for i, length := 0, len(args); i < length; i += 1 {
		idx, err := strconv.Atoi(args[i])
		if err != nil {
			listTags(issues, args)
			return
		}
		if idx > len(issues) || idx < 1 {
			fmt.Printf("Invalid issue number %d\n", idx)
			continue
		}
		if err == nil {
			b.LoadBug(bugs.Directory(bugs.GetIssuesDir() + bugs.Directory(issues[idx-1].Name())))
			b.ViewBug()
			if i < length-1 {
				fmt.Printf("\n--\n\n")
			}
		}
	}
	fmt.Printf("\n")
}

func (a BugApplication) Edit(args ArgumentList) {
	issues, _ := ioutil.ReadDir(string(bugs.GetIssuesDir()))

	// No parameters, print a list of all bugs
	if len(args) == 1 {
		idx, err := strconv.Atoi(args[0])
		if idx > len(issues) || idx < 1 {
			fmt.Printf("Invalid issue number %d\n", idx)
			return
		}
		dir := bugs.Directory(bugs.GetIssuesDir()) + bugs.Directory(issues[idx-1].Name())
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
func (a BugApplication) Close(args ArgumentList) {
	issues, _ := ioutil.ReadDir(string(bugs.GetIssuesDir()))

	// No parameters, print a list of all bugs
	if len(args) == 0 {
		fmt.Printf("Usage: %s close IssueNumber\n\nMust provide an IssueNumber to close as parameter\n", os.Args[0])
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
			dir := bugs.GetIssuesDir() + bugs.Directory(issues[idx-1].Name())
			fmt.Printf("Removing %s\n", dir)
			os.RemoveAll(string(dir))
		}
	}
}

func (a BugApplication) Purge() {
	cmd := exec.Command("git", "clean", "-fd", string(bugs.GetIssuesDir()))

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func getAllTags() []string {
	bugs := bugs.GetAllBugs()

	// Put all the tags in a map, then iterate over
	// the keys so that only unique tags are included
	tagMap := make(map[string]int, 0)
	for _, bug := range bugs {
		for _, tag := range bug.Tags() {
			tagMap[string(tag)] += 1
		}
	}

	keys := make([]string, 0, len(tagMap))
	for k := range tagMap {
		keys = append(keys, k)
	}
	return keys
}
func (a BugApplication) Tag(Args ArgumentList) {
	if len(Args) < 2 {
		fmt.Printf("Usage: %s tag issuenum tagname [more tagnames]\n", os.Args[0])
		fmt.Printf("\nBoth issue number and tagname to set are required.\n")
		fmt.Printf("\nCurrently used tags in entire tree: %s\n", strings.Join(getAllTags(), ", "))
		return
	}

	b, err := bugs.LoadBugByStringIndex(Args[0])

	if err != nil {
		fmt.Printf("Could not load bug: %s\n", err.Error())
		return
	}
	for _, tag := range Args[1:] {
		b.TagBug(bugs.Tag(tag))
	}

}
func (a BugApplication) Create(Args ArgumentList) {
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

func (a BugApplication) fieldHandler(command string, args ArgumentList,
	setCallback func(bugs.Bug, string) error, retrieveCallback func(bugs.Bug) string) {
	if len(args) < 1 {
		fmt.Printf("Usage: %s %s issuenum [set %s]\n", os.Args[0], command, command)
		return
	}

	idx, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Invalid issue number. \"%s\" is not a number.\n\n", args[0])
		fmt.Printf("Usage: %s %s issuenum [set %s]\n", os.Args[0], command, command)
		return
	}

	b, err := bugs.LoadBugByIndex(idx)
	if err != nil {
		fmt.Printf("Invalid issue number %s\n", args[0])
		return
	}
	if len(args) > 1 {
		newValue := strings.Join(args[1:], " ")
		err := setCallback(*b, newValue)
		if err != nil {
			fmt.Printf("Error setting %s: %s", command, err.Error())
		}
	} else {
		val := retrieveCallback(*b)
		if val == "" {
			fmt.Printf("%s not defined\n", command)
		} else {
			fmt.Printf("%s\n", val)
		}
	}
}
func (a BugApplication) Priority(args ArgumentList) {
	a.fieldHandler("priority", args, bugs.Bug.SetPriority, bugs.Bug.Priority)
}
func (a BugApplication) Status(args ArgumentList) {
	a.fieldHandler("status", args, bugs.Bug.SetStatus, bugs.Bug.Status)
}
func (a BugApplication) Milestone(args ArgumentList) {
	a.fieldHandler("milestone", args, bugs.Bug.SetMilestone, bugs.Bug.Milestone)
}

func (a BugApplication) Pwd() {
	fmt.Printf("%s", bugs.GetIssuesDir())
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
	type FileStatus struct {
		IndexStatus      string
		WorkingDirStatus string
		Filename         string
	}
	statusOutput := func(dir bugs.Directory) []FileStatus {
		cmd := exec.Command("git", "status", "--porcelain", "-z", string(dir))
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("Could not check git status")
			return nil
		}
		fileStatusLines := strings.Split(string(output), "\000")
		var files []FileStatus
		for _, line := range fileStatusLines {
			if len(line) == 0 {
				continue
			}
			files = append(files, FileStatus{
				IndexStatus:      line[0:1],
				WorkingDirStatus: line[1:2],
				Filename:         line[2:],
			})
		}
		return files
	}
	// Before doing anything, check git status to see if
	// the index is in a state that's going to cause an
	// error
	sOutput := statusOutput(bugs.GetIssuesDir())
	for _, file := range sOutput {
		if file.IndexStatus == "D" {
			fmt.Printf("You have manually staged changes in your issue directory which will conflict with %s commit.\n", os.Args[0])
			return
		}
	}

	sOutput = statusOutput(bugs.GetRootDir())
	for _, file := range sOutput {
		if file.IndexStatus == "A" {
			fmt.Printf("You have a new file staged in your git index, which will cause conflicts with %s commit. Please either commit your changes or unstage %s.\n", os.Args[0], file.Filename)
			return
		}
	}

	cmd := exec.Command("git", "stash", "create")

	output, err := cmd.Output()

	if err != nil {
		fmt.Printf("Could not execute git stash create")
		return
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
	cmd = exec.Command("git", "add", "-A", string(bugs.GetIssuesDir()))
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
			fmt.Printf("Error resetting the git working tree\n")
			fmt.Printf("The stash which should have your changes is: %s\n", stashHash)
		}
		cmd = exec.Command("git", "stash", "apply", "--index", stashHash)
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error restoring the git working tree")
			fmt.Printf("The stash which should have your changes is: %s\n", stashHash)
			// If nothing was stashed, it's not the end of the world.
			//fmt.Printf("Could not pop from stash\n")
		}
	}
}

func (a BugApplication) Relabel(Args ArgumentList) {
	if len(Args) < 2 {
		fmt.Printf("Usage: %s relabel issuenum New Title\n", os.Args[0])
		return
	}

	b, err := bugs.LoadBugByStringIndex(Args[0])

	if err != nil {
		fmt.Printf("Could not load bug: %s\n", err.Error())
		return
	}

	currentDir, _ := b.GetDirectory()
	newDir := bugs.GetIssuesDir() + bugs.TitleToDir(strings.Join(Args[1:], " "))
	fmt.Printf("Moving %s to %s\n", currentDir, newDir)
	err = os.Rename(string(currentDir), string(newDir))
	if err != nil {
		fmt.Printf("Error moving directory\n")
	}
}

func (a BugApplication) Version() {
	fmt.Printf("%s version 0.2-dev\n", os.Args[0])

}
