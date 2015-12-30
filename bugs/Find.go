package bugs

import (
	"fmt"
	"io/ioutil"
)

type BugNotFoundError string

func (b BugNotFoundError) Error() string {
	return "Bug not found"
}
func FindBugsByTag(tags []string) []Bug {
	return []Bug{}
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
