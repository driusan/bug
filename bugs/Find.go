package bugs

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type BugNotFoundError string

func (b BugNotFoundError) Error() string {
	return string(b)
}
func FindBugsByTag(tags []string) []Bug {
	issues, _ := ioutil.ReadDir(string(GetRootDir()) + "/issues")

	var bugs []Bug
	for idx, file := range issues {
		if file.IsDir() == true {
			bug := Bug{}
			bug.LoadBug(Directory(GetRootDir() + "/issues/" + Directory(issues[idx].Name())))
			for _, tag := range tags {
				if bug.HasTag(Tag(tag)) {
					bugs = append(bugs, bug)
					break
				}
			}
		}
	}
	return bugs
	return []Bug{}
}

func LoadBugByDirectory(dir string) (*Bug, error) {
	_, err := ioutil.ReadDir(string(GetRootDir()) + "/issues/" + dir)
	if err != nil {
		return nil, BugNotFoundError("Could not find bug " + dir)
	}
	bug := Bug{}
	bug.LoadBug(GetIssuesDir() + Directory(dir))
	return &bug, nil
}
func LoadBugByHeuristic(id string) (*Bug, error) {
	issues, _ := ioutil.ReadDir(string(GetRootDir()) + "/issues")

	idx, err := strconv.Atoi(id)
	if err == nil { // && idx > 0 && idx <= len(issues) {
		return LoadBugByIndex(idx)
	}

	var candidate *Bug
	for idx, file := range issues {
		if file.IsDir() == true {
			bug := Bug{}
			bug.LoadBug(Directory(GetRootDir() + "/issues/" + Directory(issues[idx].Name())))
			if bugid := bug.Identifier(); bugid == id {
				return &bug, nil
			} else if strings.Index(bugid, id) >= 0 {
				candidate = &bug
			}

		}
	}
	if candidate != nil {
		return candidate, nil
	}
	return nil, BugNotFoundError("Could not find bug " + id)
}
func LoadBugByStringIndex(i string) (*Bug, error) {
	issues, _ := ioutil.ReadDir(string(GetRootDir()) + "/issues")

	idx, err := strconv.Atoi(i)
	if err != nil {
		return nil, BugNotFoundError("Index not a number")
	}
	if idx < 1 || idx > len(issues) {
		return nil, BugNotFoundError("Invalid Index")
	}

	b := Bug{}
	directoryString := fmt.Sprintf("%s%s%s", GetRootDir(), "/issues/", issues[idx-1].Name())
	b.LoadBug(Directory(directoryString))
	return &b, nil
}

func LoadBugByIdentifier(id string) (*Bug, error) {
	issues, _ := ioutil.ReadDir(string(GetRootDir()) + "/issues")

	for idx, file := range issues {
		if file.IsDir() == true {
			bug := Bug{}
			bug.LoadBug(Directory(GetRootDir() + "/issues/" + Directory(issues[idx].Name())))
			if bug.Identifier() == id {
				return &bug, nil
			}
		}
	}
	return nil, BugNotFoundError("No bug named " + id)
}
func LoadBugByIndex(idx int) (*Bug, error) {
	issues, _ := ioutil.ReadDir(string(GetRootDir()) + "/issues")
	if idx < 1 || idx > len(issues) {
		return nil, BugNotFoundError("Invalid bug index")
	}

	b := Bug{}
	directoryString := fmt.Sprintf("%s%s%s", GetRootDir(), "/issues/", issues[idx-1].Name())
	b.LoadBug(Directory(directoryString))
	return &b, nil
}

func GetAllBugs() []Bug {
	issues, _ := ioutil.ReadDir(string(GetRootDir()) + "/issues")

	var bugs []Bug
	for idx, file := range issues {
		if file.IsDir() == true {
			bug := Bug{}
			bug.LoadBug(Directory(GetRootDir() + "/issues/" + Directory(issues[idx].Name())))
			bugs = append(bugs, bug)
		}
	}
	return bugs
}
