package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	args := ArgumentList(os.Args)
	if githubRepo := args.GetArgument("--github", ""); githubRepo != "" {
		if strings.Count(githubRepo, "/") != 1 {
			fmt.Fprintf(os.Stderr, "Invalid GitHub repo: %s\n", githubRepo)
			os.Exit(2)
		}
		pieces := strings.Split(githubRepo, "/")
		githubImport(pieces[0], pieces[1])

	} else if args.GetArgument("--be", "") != "" {
		beImport()
	} else {
		if strings.Count(githubRepo, "/") != 1 {
			fmt.Fprintf(os.Stderr, "Usage: %s --github user/repo\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "       %s --be\n", os.Args[0])
			fmt.Fprintf(os.Stderr, `
Use this tool to import an external bug database into the local
issues/ directory.

Either "--github user/repo" is required to import GitHub issues,
from GitHub, or "--be" is required to import a local BugsEverywhere
database.
`)
			os.Exit(2)
		}
	}
}
