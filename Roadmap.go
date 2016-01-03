package main

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"sort"
)

type BugListByMilestone [](bugs.Bug)

func (a BugListByMilestone) Len() int           { return len(a) }
func (a BugListByMilestone) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BugListByMilestone) Less(i, j int) bool { return a[i].Milestone() < a[j].Milestone() }

func (a BugApplication) Roadmap() {
	issues, _ := ioutil.ReadDir(string(bugs.GetIssuesDir()))
	var bgs [](bugs.Bug)
	for idx, _ := range issues {
		b := bugs.Bug{}
		b.LoadBug(bugs.Directory(bugs.GetIssuesDir() + bugs.Directory(issues[idx].Name())))
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
				fmt.Printf("\n## No milestone set:\n")
			} else {
				fmt.Printf("\n## %s:\n", newMilestone)
			}
		}
		fmt.Printf("- %s\n", b.Title)
		milestone = newMilestone

	}
}
