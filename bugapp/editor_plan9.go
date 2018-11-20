package bugapp

import (
	"os"
)

func getEditor() string {
	editor := os.Getenv("editor")

	if editor != "" {
		return editor
	}
	return "sam"

}

