package main

import (
	"fmt"
	"os"
	//"regex"
	"github.com/driusan/bug/bugs"
)

func getEditor() string {
	editor := os.Getenv("EDITOR")

	if editor != "" {
		return editor
	}
	return "vim"

}

func main() {
	app := BugApplication{}
	if bugs.GetRootDir() == "" {
		fmt.Printf("Could not find issues directory.\n")
		fmt.Printf("Make sure either the PMIT environment variable is set, or a parent directory of your working directory has an issues folder.\n")
		fmt.Printf("Aborting.\n")
		os.Exit(2)
	}
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "add":
			fallthrough
		case "new":
			fallthrough
		case "create":
			app.Create(os.Args[2:])
		case "view":
			fallthrough
		case "list":
			app.List(os.Args[2:])
		case "priority":
			app.Priority(os.Args[2:])
		case "status":
			app.Status(os.Args[2:])
		case "milestone":
			app.Milestone(os.Args[2:])
		case "tag":
			app.Tag(os.Args[2:])
		case "mv":
			fallthrough
		case "rename":
			fallthrough
		case "retitle":
			fallthrough
		case "relabel":
			app.Relabel(os.Args[2:])
		case "purge":
			app.Purge()
		case "rm":
			fallthrough
		case "close":
			app.Close(os.Args[2:])
		case "edit":
			app.Edit(os.Args[2:])
		case "env":
			app.Env()
		case "dir":
			fallthrough
		case "pwd":
			app.Pwd()
		case "commit":
			app.Commit()
		case "roadmap":
			app.Roadmap()
		case "help":
			fallthrough
		default:
			app.Help(os.Args[1:]...)
		}
	} else {
		app.Help()
	}
}
