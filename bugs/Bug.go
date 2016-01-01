package bugs

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type Bug struct {
	Title       string
	Description string
}

type Status string
type Priority string
type Milestone string
type Tag string

func TitleToDir(title string) Directory {
	re := regexp.MustCompile("(-+)")
	s := re.ReplaceAllString(title, "-$1")
	s = strings.Replace(s, " ", "-", -1)
    return Directory(s)
}
func (b Bug) GetDirectory() (Directory, error) {
	return GetRootDir() + "/issues/" + TitleToDir(b.Title), nil
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

func (b *Bug) TagBug(tag Tag) {
	if dir, err := b.GetDirectory(); err == nil {
		os.Mkdir(string(dir)+"/tags/", 0755)
		ioutil.WriteFile(string(dir)+"/tags/"+string(tag), []byte(""), 0644)
	} else {
		fmt.Printf("Error tagging bug: %s", err.Error())
	}
}
func (b Bug) ViewBug() {
	fmt.Printf("Title: %s\n\n", b.Title)
	fmt.Printf("Description:\n%s", b.Description)

	status := b.Status()
	if status != "" {
		fmt.Printf("\nStatus: %s", status)
	}
	priority := b.Priority()
	if priority != "" {
		fmt.Printf("\nPriority: %s", priority)
	}
	milestone := b.Milestone()
	if milestone != "" {
		fmt.Printf("\nMilestone: %s", milestone)
	}
	tags := b.StringTags()
	if tags != nil {
		fmt.Printf("\nTags: %s", strings.Join([]string(tags), ", "))
	}

}

func (b Bug) StringTags() []string {
	dir, _ := b.GetDirectory()
	dir += "/tags/"
	issues, err := ioutil.ReadDir(string(dir))
	if err != nil {
		return nil
	}

	tags := make([]string, 0, len(issues))
	for _, issue := range issues {
		tags = append(tags, issue.Name())
	}
	return tags
}
func (b Bug) Tags() []Tag {
	dir, _ := b.GetDirectory()
	dir += "/tags/"
	issues, err := ioutil.ReadDir(string(dir))
	if err != nil {
		return nil
	}

	tags := make([]Tag, 0, len(issues))
	for _, issue := range issues {
		tags = append(tags, Tag(issue.Name()))
	}
	return tags

}

func (b Bug) getField(fieldName string) string {
	dir, _ := b.GetDirectory()
	field, err := ioutil.ReadFile(string(dir) + "/" + fieldName)
	if err != nil {
		return ""
	}
	lines := strings.Split(string(field), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}
	return ""
}

func (b Bug) setField(fieldName, value string) error {
	dir, _ := b.GetDirectory()
	oldValue, err := ioutil.ReadFile(string(dir) + "/" + fieldName)
	var oldLines []string
	if err == nil {
		oldLines = strings.Split(string(oldValue), "\n")
	}

	newValue := ""
	if len(oldLines) >= 1 {
		// If there were 0 or 1 old lines, overwrite them
		oldLines[0] = value
		newValue = strings.Join(oldLines, "\n")
	} else {
		newValue = value
	}

	err = ioutil.WriteFile(string(dir)+"/"+fieldName, []byte(newValue), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (b Bug) Status() string {
	return b.getField("Status")
}

func (b Bug) SetStatus(newStatus string) error {
	return b.setField("Status", newStatus)
}
func (b Bug) Priority() string {
	return b.getField("Priority")
}

func (b Bug) SetPriority(newValue string) error {
	return b.setField("Priority", newValue)
}
func (b Bug) Milestone() string {
	return b.getField("Milestone")
}

func (b Bug) SetMilestone(newValue string) error {
	return b.setField("Milestone", newValue)
}
