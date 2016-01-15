package main

import (
	"encoding/json"
	"fmt"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

func beImportComments(b bugs.Bug, directory string, includeHeaders bool) string {
	file := directory + "/body"
	data, _ := ioutil.ReadFile(file)
	if includeHeaders == false {
		return string(data)
	}

	/* Be appears to store comment metadata in a JSON file named values
	   with the format:
	   {
	       "Author": "Dave MacFarlane <driusan@gmail.com>",
	       "Content-type": "text/plain",
	       "Date": "Tue, 12 Jan 2016 21:44:24 +0000"
	   }
	*/
	type BeValues struct {
		Author      string `json:Author`
		ContentType string `json:Content-type`
		Date        string `json:Date`
	}
	file = directory + "/values"
	jsonVal, _ := ioutil.ReadFile(file)
	var beComment BeValues
	json.Unmarshal([]byte(jsonVal), &beComment)
	return "---------- Comment ---------\nFrom:" + beComment.Author + "\nDate:" + beComment.Date + "\n\n" + string(data)
}

func beImportBug(identifier, issuesDir, fullbepath string) {
	/* BE appears to store the top level data of a bug
	   ins a json file named values with the format:
	    {
	        "creator": "Dave MacFarlane <driusan@gmail.com>",
	        "reporter": "Dave MacFarlane <driusan@gmail.com>",
	        "severity": "minor",
	        "status": "open",
	        "summary": "abc",
	        "time": "Tue, 12 Jan 2016 00:05:28 +0000"
	    }

	    and the description of bugs entirely in comments.
	    All we really care about is the summary so that we
	    can get the directory name for the issues/ directory,
	    but the severity+status can also be used as a status
	    to ensure that we have at least 1 file to be tracked
	    by git.
	*/

	type BeValues struct {
		Creator  string `json:creator`
		Reporter string `json:reporter`
		Severity string `json:severity`
		Status   string `json:status`
		Summary  string `json:summary`
		Time     string `json:time`
	}
	file := fullbepath + "/values"

	fmt.Printf("File: %s\n", file)
	data, _ := ioutil.ReadFile(file)
	var beBug BeValues
	err := json.Unmarshal([]byte(data), &beBug)
	if err != nil {
		fmt.Printf("Error unmarshalling data: %s\n", err.Error())
	}

	fmt.Printf("%s\n", beBug)

	bugdir := bugs.TitleToDir(beBug.Summary)

	b := bugs.Bug{bugs.Directory(issuesDir) + bugdir}
	if dir := b.GetDirectory(); dir != "" {
		os.Mkdir(string(dir), 0755)
	}
	if beBug.Status != "" && beBug.Severity != "" {
		b.SetStatus(beBug.Status + ":" + beBug.Severity)
	}

	comments := fullbepath + "/comments/"
	dir, err := os.Open(comments)

	files, err := dir.Readdir(-1)
	var Description string
	if len(files) > 0 && err == nil {
		for _, file := range files {
			if file.IsDir() {
				Description = Description + "\n" +
					beImportComments(b, comments+file.Name(), len(files) > 1)
			}
		}
	}
	b.SetDescription(Description)
	b.SetIdentifier(identifier)
}
func beImportBugs(prefix, issuesdir, bedir, dirname string) {
	bugsdir := bedir + "/" + dirname + "/bugs"
	dir, err := os.Open(bugsdir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open bug directory %s\n", bugsdir)
	}
	files, _ := dir.Readdir(-1)

	lastIdentifier := ""
	nextIdentifier := ""
	sort.Sort(fileSorter(files))
	for idx, file := range files {
		if file.IsDir() {
			if idx < len(files)-1 {
				nextIdentifier = files[idx+1].Name()
			}
			name := shortestPrefix(file.Name(), nextIdentifier, lastIdentifier, 3)
			identifier := fmt.Sprintf("%s/%s", prefix, name)
			beImportBug(identifier, issuesdir, bugsdir+"/"+file.Name())
			lastIdentifier = file.Name()
		}
	}
}

type fileSorter []os.FileInfo

func (a fileSorter) Len() int {
	return len(a)
}
func (a fileSorter) Less(i, j int) bool {
	return a[i].Name() < a[j].Name()
}
func (a fileSorter) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func beImport() {
	wd, err := os.Getwd()
	if dir := walkAndSearch(wd, []string{".be"}); err != nil || dir == nil {
		fmt.Fprintf(os.Stderr, "Could not find any Bugs Everywhere repository relative to current path.\n")
		os.Exit(3)
	} else {
		files, err := dir.Readdir(-1)
		sort.Sort(fileSorter(files))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error importing BE database: %s\n", err.Error())
			os.Exit(4)
		}

		issuesDir := bugs.GetIssuesDir()
		lastIdentifier := ""
		nextIdentifier := ""
		for idx, file := range files {
			if file.IsDir() {
				if idx < len(files)-1 {
					nextIdentifier = files[idx+1].Name()
				}
				name := shortestPrefix(file.Name(), nextIdentifier, lastIdentifier, 3)

				beImportBugs(name, string(issuesDir), dir.Name(), file.Name())
				lastIdentifier = file.Name()
			}
		}
	}
}

func walkAndSearch(startpath string, dirnames []string) *os.File {
	for _, dirname := range dirnames {
		if dirinfo, err := os.Stat(startpath + "/" + dirname); err == nil && dirinfo.IsDir() {
			file, err := os.Open(startpath + "/" + dirname)
			if err != nil {
				return nil
			}
			return file
		}
	}

	pieces := strings.Split(startpath, "/")

	for i := len(pieces); i > 0; i -= 1 {
		dir := strings.Join(pieces[0:i], "/")
		for _, dirname := range dirnames {
			if dirinfo, err := os.Stat(dir + "/" + dirname); err == nil && dirinfo.IsDir() {
				file, err := os.Open(dir + "/" + dirname)
				if err != nil {
					return nil
				}
				return file
			}
		}
	}
	return nil
}

func shortestPrefix(name, nextIdentifier, lastIdentifier string, min int) string {
	for i := 0; i < len(name); i += 1 {
		if i < min {
			continue
		}
		if nextIdentifier != "" && len(nextIdentifier) <= i {
			if name[0:i] == nextIdentifier[0:i] {
				continue
			}
		}

		if lastIdentifier != "" && len(lastIdentifier) <= i {
			if name[0:i] == lastIdentifier[0:i] {
				continue
			}
		}
		return name[0:i]
	}
	return name
}
