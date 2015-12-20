package bugs

import (
	"fmt"
	"io/ioutil"
	"os"
	//"regex"
	"strings"
)

type Bug struct {
	Title       string
	Description string
}

func (b Bug) GetDirectory() (Directory, error) {
	//re := regexp.MustCompile("-[-*]
	s := strings.Replace(b.Title, "-", "--", -1)

	tokens := strings.Split(s, " ")

	return GetRootDir() + "/issues/" + Directory(strings.Join(tokens, "-")), nil
}

func (b *Bug) LoadBug(dir Directory) {
	b.Title = dir.GetShortName().ToTitle()

	desc, err := ioutil.ReadFile(string(dir) + "/Description")

	if err != nil {
		b.Description = "No description provided"
		return
	}

	b.Description = string(desc)
}

func (b *Bug) TagBug(tag string) {
	if dir, err := b.GetDirectory(); err == nil {
		os.Mkdir(string(dir)+"/tags/", 0755)
		ioutil.WriteFile(string(dir)+"/tags/"+tag, []byte(""), 0644)
	} else {
		fmt.Printf("Error tagging bug: %s", err.Error())
	}
}
func (b Bug) ViewBug() {
	fmt.Printf("Title: %s\n\n", b.Title)
	fmt.Printf("Description:\n%s", b.Description)

	tags := b.Tags()
	if tags != nil {
		fmt.Printf("\nTags: %s", strings.Join(tags, ", "))
	}
}

func (b Bug) Tags() []string {
	dir, _ := b.GetDirectory()
	dir += "/tags/"
	issues, err := ioutil.ReadDir(string(dir))
	if err != nil {
		return nil
	}

	tags := []string{}
	for _, issue := range issues {
		tags = append(tags, issue.Name())
	}
	return tags

}
