package bugs

import (
	"fmt"
	"io/ioutil"
	"strconv"
)

type BugNotFoundError string

func (b BugNotFoundError) Error() string {
	return string(b)
}
func FindBugsByTag(tags []string) []Bug {
	return []Bug{}
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

func LoadBugByIndex(idx int) (*Bug, error) {
	issues, _ := ioutil.ReadDir(string(GetRootDir()) + "/issues")
	if idx < 1 || idx > len(issues) {
		return nil, BugNotFoundError("Invalid Index")
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
