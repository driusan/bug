// bug writes code problem reports to plain text files.
package main

import (
	"fmt"
	"github.com/driusan/bug/bugapp"
	"github.com/driusan/bug/bugs"
	"os"
)

func main() {
	var skipRootCheck bool = false
	switch len(os.Args) {
	case 0, 1:
		skipRootCheck = true
	case 2:
		if os.Args[1] == "help" {
			skipRootCheck = true
		}
	case 3:
		if os.Args[2] == "--help" {
			skipRootCheck = true
		}

	}
	if skipRootCheck == false && bugs.GetRootDir() == "" {
		fmt.Printf("Could not find issues directory.\n")
		fmt.Printf("Make sure either the PMIT environment variable is set, or a parent directory of your working directory has an issues folder.\n")
		fmt.Println("(If you just started new repo, you probably want to create directory named `issues`).")
		fmt.Printf("Aborting.\n")
		os.Exit(2)

	}

	if len(os.Args) > 1 {
		if len(os.Args) >= 3 && os.Args[2] == "--help" {
			os.Args[1], os.Args[2] = "help", os.Args[1]
		}
		switch os.Args[1] {
		case "add", "new", "create":
			bugapp.Create(os.Args[2:])
		case "view", "list":
			// bug list with no parameters shouldn't autopage,
			// bug list with bugs to view should. So the original
			// stdout is passed as a parameter.
			bugapp.List(os.Args[2:])
		case "priority":
			bugapp.Priority(os.Args[2:])
		case "status":
			bugapp.Status(os.Args[2:])
		case "milestone":
			bugapp.Milestone(os.Args[2:])
		case "id", "identifier":
			bugapp.Identifier(os.Args[2:])
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
			bugapp.Commit(os.Args[2:])
		case "roadmap":
			bugapp.Roadmap(os.Args[2:])
		case "find":
			bugapp.Find(os.Args[2:])
		case "help":
			fallthrough
		default:
			bugapp.Help(os.Args[1:]...)
		}
	} else {
		bugapp.Help()
	}
}
