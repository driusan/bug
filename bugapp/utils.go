package bugapp

import (
    "os"
)

func getEditor() string {
	editor := os.Getenv("EDITOR")

	if editor != "" {
		return editor
	}
	return "vim"

}

type ArgumentList []string

func (args ArgumentList) HasArgument(arg string) bool {
	for _, argCandidate := range args {
		if arg == argCandidate {
			return true
		}
	}
	return false
}
