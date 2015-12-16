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

func getRootDir() string {
	dir := os.Getenv("PMIT")
	if dir != "" {
		return dir
	}

	wd, _ := os.Getwd()

	if dirinfo, err := os.Stat(wd + "/issues"); err == nil && dirinfo.IsDir() {
		return wd
	}

	// There's no environment variable and no issues
	// directory, so walk up the tree until we find one
	pieces := strings.Split(wd, "/")

	for i := len(pieces); i > 0; i -= 1 {
		dir := strings.Join(pieces[0:i], "/")
		if dirinfo, err := os.Stat(dir + "/issues"); err == nil && dirinfo.IsDir() {
			return dir
		}
	}
	return ""
}

type Directory string

func (d Directory) GetShortName() Directory {
	pieces := strings.Split(string(d), "/")
	return Directory(pieces[len(pieces)-1])
}

func (d Directory) ToTitle() string {
	tokens := strings.Split(string(d), "-")
	return strings.Join(tokens, " ")
}

type Bug struct {
	Title       string
	Description string
}

func (b Bug) GetDirectory() (Directory, error) {
	//re := regexp.MustCompile("-[-*]
	s := strings.Replace(b.Title, "-", "--", -1)

	tokens := strings.Split(s, " ")

	return Directory(getRootDir() + "/issues/" + strings.Join(tokens, "-")), nil
}

func (b *Bug) LoadBug(dir Directory) {
	b.Title = dir.GetShortName().ToTitle()

	desc, err := ioutil.ReadFile(string(dir) + "/Description")

	if err != nil {
		b.Description = "No description provided"
		return
	}

	b.Description = string(desc)
}

func (b Bug) ViewBug() {
	fmt.Printf("\nTitle: %s\n\n", b.Title)
	fmt.Printf("Description:\n%s\n", b.Description)
}
func printHelp() {
	fmt.Printf("Usage: " + os.Args[0] + " command [options]\n\n")
	fmt.Printf("Valid commands\n")
	fmt.Printf("\tcreate\tFile a new bug\n")
	fmt.Printf("\tlist\tList existing bugs\n")
	fmt.Printf("\tclose\tDelete an existing bug\n")
	fmt.Printf("\trm\tAlias of close\n")
	fmt.Printf("\tenv\tShow settings that bug will use if invoked from this directory\n")
	fmt.Printf("\tdir\tPrints the issues directory to stdout (useful subcommand in the shell)\n")
	fmt.Printf("\thelp\tShow this screen\n")
}

func showEnv() {
	fmt.Printf("Settings used by this command:\n")
	fmt.Printf("\nIssues directory:\t%s/issues/", getRootDir())
	fmt.Printf("\nEditor:\t%s", getEditor())
	fmt.Printf("\n")
}
func listBugs(args []string) {
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
func closeBugs(args []string) {
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

func getEditor() string {
	editor := os.Getenv("EDITOR")

	if editor != "" {
		return editor
	}
	return "vim"

}
func createBug(Args []string) {
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

func printDir() {
	fmt.Printf("%s", getRootDir()+"/issues")
}
func main() {
	if getRootDir() == "" {
		fmt.Printf("Could not find issues directory.\n")
		fmt.Printf("Make sure either the PMIT environment variable is set, or a parent directory of your working directory has an issues folder.\n")
		fmt.Printf("Aborting.\n")
		os.Exit(2)
	}
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "create":
			createBug(os.Args[2:])
		case "view":
			fallthrough
		case "list":
			listBugs(os.Args[2:])
		case "rm":
			fallthrough
		case "close":
			closeBugs(os.Args[2:])
		case "env":
			showEnv()
		case "dir":
			printDir()
		case "help":
			fallthrough
		default:
			printHelp()
		}
	} else {
		printHelp()
	}
}
