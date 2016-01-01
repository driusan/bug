package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	//"regex"
	"github.com/driusan/bug/bugs"
	"strings"
)

type BugApplication struct{}

func (a BugApplication) Env() {
	fmt.Printf("Settings used by this command:\n")
	fmt.Printf("\nIssues directory:\t%s/issues/", bugs.GetRootDir())
	fmt.Printf("\nEditor:\t%s", getEditor())
	fmt.Printf("\n")
}

func listTags(files []os.FileInfo, args []string) {
	hasTag := func(tags []string, tag string) bool {
		for _, candidate := range tags {
			if candidate == tag {
				return true
			}
		}
		return false
	}
	b := bugs.Bug{}
	for idx, _ := range files {
		b.LoadBug(bugs.Directory(bugs.GetRootDir() + "/issues/" + bugs.Directory(files[idx].Name())))

		tags := b.Tags()
		for _, tag := range args {
			if hasTag(tags, tag) {
				fmt.Printf("Issue %d: %s (%s)\n", idx+1, b.Title, strings.Join(tags, ", "))
			}
		}
	}
}
func (a BugApplication) List(args []string) {
	issues, _ := ioutil.ReadDir(string(bugs.GetRootDir()) + "/issues")

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
			b.LoadBug(bugs.Directory(bugs.GetRootDir() + "/issues/" + bugs.Directory(issues[idx-1].Name())))
			b.ViewBug()
			if i < length-1 {
				fmt.Printf("\n--\n\n")
			}
		}
	}
	fmt.Printf("\n")
}

func (a BugApplication) Edit(args []string) {
	issues, _ := ioutil.ReadDir(string(bugs.GetRootDir()) + "/issues")

	// No parameters, print a list of all bugs
	if len(args) == 1 {
		idx, err := strconv.Atoi(args[0])
		if idx > len(issues) || idx < 1 {
			fmt.Printf("Invalid issue number %d\n", idx)
			return
		}
		dir := bugs.Directory(bugs.GetRootDir()) + "/issues/" + bugs.Directory(issues[idx-1].Name())
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
	issues, _ := ioutil.ReadDir(string(bugs.GetRootDir()) + "/issues")

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
			dir := bugs.GetRootDir() + "/issues/" + bugs.Directory(issues[idx-1].Name())
			fmt.Printf("Removing %s\n", dir)
			os.RemoveAll(string(dir))
		}
	}
}

func (a BugApplication) Purge() {
	cmd := exec.Command("git", "clean", "-fd", string(bugs.GetRootDir())+"/issues")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func getAllTags() []string {
	issues, _ := ioutil.ReadDir(string(bugs.GetRootDir()) + "/issues")
	tagMap := make(map[string]bool)
	for idx, _ := range issues {
		var b bugs.Bug
		b.LoadBug(bugs.Directory(bugs.GetRootDir() + "/issues/" + bugs.Directory(issues[idx].Name())))
		for _, tag := range b.Tags() {
			tagMap[tag] = true
		}
	}

	keys := make([]string, 0, len(tagMap))
	for k := range tagMap {
		keys = append(keys, k)
	}
	return keys
}
func (a BugApplication) Tag(Args []string) {
	if len(Args) < 2 {
		fmt.Printf("Usage: %s tag issuenum tagname [more tagnames]\n", os.Args[0])
		fmt.Printf("\nBoth issue number and tagname to set are required.\n")
		fmt.Printf("\nCurrently used tags in entire tree: %s\n", strings.Join(getAllTags(), ", "))
		return
	}

	issues, err := ioutil.ReadDir(string(bugs.GetRootDir()) + "/issues")
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
	var b bugs.Bug
	b.LoadBug(bugs.Directory(bugs.GetRootDir() + "/issues/" + bugs.Directory(issues[idx].Name())))
	for _, tag := range Args[1:] {
		b.TagBug(tag)
	}

}
func (a BugApplication) Create(Args []string) {
	var noDesc bool = false

	if Args != nil && Args[0] == "-n" {
		noDesc = true
		Args = Args[1:]
	}

	var bug bugs.Bug
	bug = bugs.Bug{
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

func (a BugApplication) fieldHandler(command string, args []string,
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
func (a BugApplication) Priority(args []string) {
	a.fieldHandler("priority", args, bugs.Bug.SetPriority, bugs.Bug.Priority)
}
func (a BugApplication) Status(args []string) {
	a.fieldHandler("status", args, bugs.Bug.SetStatus, bugs.Bug.Status)
}
func (a BugApplication) Milestone(args []string) {
	a.fieldHandler("milestone", args, bugs.Bug.SetMilestone, bugs.Bug.Milestone)
}

func (a BugApplication) Dir() {
	fmt.Printf("%s", bugs.GetRootDir()+"/issues")
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
	cmd = exec.Command("git", "add", "-A", string(bugs.GetRootDir())+"/issues")
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

type BugListByMilestone [](bugs.Bug)

func (a BugListByMilestone) Len() int           { return len(a) }
func (a BugListByMilestone) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BugListByMilestone) Less(i, j int) bool { return a[i].Milestone() < a[j].Milestone() }

func (a BugApplication) Roadmap() {
	issues, _ := ioutil.ReadDir(string(bugs.GetRootDir()) + "/issues")
	var bgs [](bugs.Bug)
	for idx, _ := range issues {
		b := bugs.Bug{}
		b.LoadBug(bugs.Directory(bugs.GetRootDir() + "/issues/" + bugs.Directory(issues[idx].Name())))
		bgs = append(bgs, b)
	}
	sort.Sort(BugListByMilestone(bgs))

	fmt.Printf("# Roadmap for %s\n", bugs.GetRootDir().GetShortName().ToTitle())
	milestone := ""
	for i := len(bgs) - 1; i >= 0; i -= 1 {
		b := bgs[i]
		newMilestone := b.Milestone()
		if milestone != newMilestone {
			if newMilestone == "" {
				fmt.Printf("\nNo milestone set:\n")
			} else {
				fmt.Printf("\n%s:\n", newMilestone)
			}
		}
		fmt.Printf("- %s\n", b.Title)
		milestone = newMilestone

	}
}
