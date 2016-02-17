package scm

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type GitManager struct {
	Autoclose bool
}

func (a GitManager) Purge(dir bugs.Directory) error {
	cmd := exec.Command("git", "clean", "-fd", string(dir))

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (a GitManager) closedGithubIssues(dir bugs.Directory) []string {
	ghRegex := regexp.MustCompile("(?im)^-Github:(.*)$")
	cmd := exec.Command("git", "status", "-z", "--porcelain", string(dir))
	out, _ := cmd.CombinedOutput()
	files := strings.Split(string(out), "\000")
	var retVal []string
	for _, file := range files {
		if file == "" {
			continue
		}

		if file[:1] == "D" && strings.HasSuffix(file, "Identifier") {
			diff := exec.Command("git", "diff", "--staged", "--", file[3:])
			diffout, _ := diff.CombinedOutput()
			//fmt.Printf("Output: %s", diffout)
			if matches := ghRegex.FindStringSubmatch(string(diffout)); len(matches) > 1 {
				retVal = append(retVal, strings.TrimSpace(matches[1]))
			}
		}
	}
	return retVal
}

func (a GitManager) Commit(dir bugs.Directory, commitMsg string) error {
	cmd := exec.Command("git", "add", "-A", string(dir))
	if err := cmd.Run(); err != nil {
		fmt.Printf("Could not add issues to be commited: %s?\n", err.Error())
		return err

	}

	var closesGH string
	if a.Autoclose {
		ci := a.closedGithubIssues(dir)
		if len(ci) > 0 {
			closesGH = fmt.Sprintf("\nCloses %s\n", strings.Join(ci, ", closes "))
		}
	}

	file, err := ioutil.TempFile("", "bugCommit")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create temporary file.\nNothing commited.\n")
		return err
	}
	defer os.Remove(file.Name())

	fmt.Fprintf(file, "%s\n%s", commitMsg, closesGH)
	cmd = exec.Command("git", "commit", "-o", string(dir), "-F", file.Name(), "-q")
	if err := cmd.Run(); err != nil {
		// If nothing was added commit will have an error,
		// but we don't care it just means there's nothing
		// to commit.
		fmt.Printf("No new issues commited\n")
		return nil
	}
	return nil
}

func (a GitManager) GetSCMType() string {
	return "git"
}
