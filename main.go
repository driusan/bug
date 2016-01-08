package main

import (
	"fmt"
	"os"
	"github.com/driusan/bug/bugapp"
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
	if bugs.GetRootDir() == "" {
		fmt.Printf("Could not find issues directory.\n")
		fmt.Printf("Make sure either the PMIT environment variable is set, or a parent directory of your working directory has an issues folder.\n")
		fmt.Printf("Aborting.\n")
		os.Exit(2)
	}
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "add", "new", "create":
			bugapp.Create(os.Args[2:])
		case "view", "list":
			bugapp.List(os.Args[2:])
		case "priority":
			bugapp.Priority(os.Args[2:])
		case "status":
			bugapp.Status(os.Args[2:])
		case "milestone":
			bugapp.Milestone(os.Args[2:])
		case "tag":
			bugapp.Tag(os.Args[2:])
		case "mv", "rename", "retitle", "relabel":
			bugapp.Relabel(os.Args[2:])
		case "purge":
			bugapp.Purge()
		case "rm", "close":
			bugapp.Close(os.Args[2:])
		case "edit":
			bugapp.Edit(os.Args[2:])
		case "--version", "version":
			bugapp.Version()
		case "env":
			bugapp.Env()
		case "dir", "pwd":
			bugapp.Pwd()
		case "commit":
			bugapp.Commit()
		case "roadmap":
			bugapp.Roadmap(os.Args[2:])
		case "help":
			fallthrough
		default:
			bugapp.Help(os.Args[1:]...)
		}
	} else {
		bugapp.Help()
	}
}
