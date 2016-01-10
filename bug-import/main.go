package main

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"github.com/google/go-github/github"
	"os"
	"strings"
)

func githubImport(user, repo string) {
	client := github.NewClient(nil)
	issueDir := bugs.GetIssuesDir()
	issues, _, err := client.Issues.ListByRepo(user, repo, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for _, issue := range issues {
		if issue.PullRequestLinks == nil {
			b := bugs.Bug{issueDir + bugs.TitleToDir(*issue.Title)}
			if dir := b.GetDirectory(); dir != "" {
				os.Mkdir(string(dir), 0755)
			}
			if issue.Body != nil {
				b.SetDescription(*issue.Body)
			}
			if issue.Milestone != nil {
				b.SetMilestone(*issue.Milestone.Title)
			}
			// Don't set a bug identifier, but put an empty line and
			// then a GitHub identifier, so that bug commit can include
			// "Closes ..." in the commit message.
			b.SetIdentifier(fmt.Sprintf("\n\nGitHub:%s/%s%s%d\n", user, repo, "#", *issue.Number))
			for _, l := range issue.Labels {
				b.TagBug(bugs.Tag(*l.Name))
			}
			fmt.Printf("Importing %s\n", *issue.Title)
		}
	}
}
func main() {

	args := ArgumentList(os.Args)
	if githubRepo := args.GetArgument("--github", ""); githubRepo != "" {
		if strings.Count(githubRepo, "/") != 1 {
			fmt.Fprintf(os.Stderr, "Invalid GitHub repo: %s\n", githubRepo)
			os.Exit(2)
		}
		pieces := strings.Split(githubRepo, "/")
		githubImport(pieces[0], pieces[1])

	}
}
