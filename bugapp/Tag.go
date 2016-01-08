package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"os"
	"strings"
)

func getAllTags() []string {
	bugs := bugs.GetAllBugs()

	// Put all the tags in a map, then iterate over
	// the keys so that only unique tags are included
	tagMap := make(map[string]int, 0)
	for _, bug := range bugs {
		for _, tag := range bug.Tags() {
			tagMap[string(tag)] += 1
		}
	}

	keys := make([]string, 0, len(tagMap))
	for k := range tagMap {
		keys = append(keys, k)
	}
	return keys
}
func Tag(Args ArgumentList) {
	if len(Args) < 2 {
		fmt.Printf("Usage: %s tag [--rm] issuenum tagname [more tagnames]\n", os.Args[0])
		fmt.Printf("\nBoth issue number and tagname to set are required.\n")
		fmt.Printf("\nCurrently used tags in entire tree: %s\n", strings.Join(getAllTags(), ", "))
		return
	}
	var removeTags bool = false
	if Args[0] == "--rm" {
		removeTags = true
		Args = Args[1:]
	}

	b, err := bugs.LoadBugByStringIndex(Args[0])

	if err != nil {
		fmt.Printf("Could not load bug: %s\n", err.Error())
		return
	}
	for _, tag := range Args[1:] {
		if removeTags {
			b.RemoveTag(bugs.Tag(tag))
		} else {
			b.TagBug(bugs.Tag(tag))
		}
	}

}
