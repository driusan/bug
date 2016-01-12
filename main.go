package main

import (
	"fmt"
	"github.com/driusan/bug/bugapp"
	"github.com/driusan/bug/bugs"
	"os"
	"os/exec"
	//    "bytes"
	//   "io"
)

func main() {
	if bugs.GetRootDir() == "" {
		fmt.Printf("Could not find issues directory.\n")
		fmt.Printf("Make sure either the PMIT environment variable is set, or a parent directory of your working directory has an issues folder.\n")
		fmt.Printf("Aborting.\n")
		os.Exit(2)
	}

	// Create a pipe for a pager to use
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	// Capture STDOUT for the Pager
	stdout := os.Stdout
	os.Stdout = w

	// Invoke less -RF attached to the pipe
	// we created
	cmd := exec.Command("less", "-RF")
	cmd.Stdin = r
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr
	// Make sure the pipe is closed after we
	// finish, then restore STDOUT
	defer func() {
		w.Close()
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Stdout = stdout
	}()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "add", "new", "create":
			os.Stdout = stdout
			bugapp.Create(os.Args[2:])
		case "view", "list":
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
			os.Stdout = stdout
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
		case "help":
			fallthrough
		default:
			bugapp.Help(os.Args[1:]...)
		}
	} else {
		bugapp.Help()
	}
}
