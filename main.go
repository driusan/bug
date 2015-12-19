package main

import (
	"fmt"
	"os"
	//"regex"
	"strings"
	"github.com/driusan/bug/bugs"
)

func getEditor() string {
	editor := os.Getenv("EDITOR")

	if editor != "" {
		return editor
	}
	return "vim"

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
		case "create":
			app.Create(os.Args[2:])
		case "view":
			fallthrough
		case "list":
			app.List(os.Args[2:])
		case "tag":
			app.Tag(os.Args[2:])
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
		case "pwd":
			fallthrough
		case "dir":
			app.Dir()
		case "commit":
			app.Commit()
		case "help":
			fallthrough
		default:
			app.Help(os.Args[1:]...)
		}
	} else {
		app.Help()
	}
}
