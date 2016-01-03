package main

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"sort"
    "github.com/blang/semver"
    "strconv"
)

type BugListByMilestone [](bugs.Bug)

func (a BugListByMilestone) Len() int           { return len(a) }
func (a BugListByMilestone) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BugListByMilestone) Less(i, j int) bool {
    iMS := a[i].Milestone()
    jMS := a[j].Milestone()
    // If there's a "v" at the start, strip it out
    // before doing any comparisons of semantic
    // versions
    if iMS[0] == "v"[0] {
        iMS = iMS[1:]
    }
    if jMS[0] == "v"[0] {
        jMS = jMS[1:]
    }
    // First try semantic versioning comparison
    iVer, iVerErr := semver.Make(iMS)
    jVer, jVerErr := semver.Make(jMS)
    if iVerErr == nil && jVerErr == nil {
        return iVer.LT(jVer)
    }

    // Next try floating point comparison as an 
    // approximation of real number comparison..
    iFloat, iVerErr :=strconv.ParseFloat(iMS, 32)
    jFloat, jVerErr :=strconv.ParseFloat(jMS, 32)
    if iVerErr == nil && jVerErr == nil {
        return iFloat < jFloat
    }

    // Finally, just use a normal string collation
    return iMS < jMS
}

func (a BugApplication) Roadmap() {
    bgs := bugs.GetAllBugs()
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
