package main

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"github.com/google/go-github/github"
	"os"
	"context"
)

func githubImport(user, repo string) {
	client := github.NewClient(nil)
	issueDir := bugs.GetIssuesDir()
	opt := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	issues, resp, err := client.Issues.ListByRepo(context.Background(), user, repo, opt)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for lastPage := false; lastPage != true; {
		for _, issue := range issues {
			if issue.PullRequestLinks == nil {
				b := bugs.Bug{Dir: issueDir + bugs.TitleToDir(*issue.Title)}
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
		if resp.NextPage == 0 {
			lastPage = true
		} else {
			opt.ListOptions.Page = resp.NextPage
			issues, resp, err = client.Issues.ListByRepo(context.Background(), user, repo, opt)
		}
	}
}
