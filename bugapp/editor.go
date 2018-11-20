// +build !plan9

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

