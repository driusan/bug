package main

import (
	"fmt"
	"github.com/driusan/GoWebapp/URLHandler"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type MainPageHandler struct {
	URLHandler.DefaultHandler
}

type BugPageHandler struct {
	URLHandler.DefaultHandler
}
func (m MainPageHandler) Get(*http.Request, map[string]interface{}) (string, error) {
	issues, _ := ioutil.ReadDir(bugs.GetRootDir() + "/issues")

	var ret string = "<html><body><h2>Open issues</h2><ol>"
	for _, issue := range issues {
		var dir bugs.Directory = bugs.Directory(issue.Name())
		ret += fmt.Sprintf("<li><a href=\"/issues/%s\">%s</a></li>\n", (dir), dir.ToTitle())
	}

	ret += "</ol></body></html>"

	return ret, nil
}
func (m BugPageHandler) Get(r *http.Request, extras map[string]interface{}) (string, error) {
	bugDir := getRootDir() + r.URL.Path
	b := bugs.Bug{}
	b.LoadBug(bugs.Directory(bugDir))
	return "<html><body><h2>" + b.Title + "</h2>" + "<p>" + strings.Replace(b.Description, "\n\n", "</p><p>", -1) + "</p></body></html>", nil
}

func main() {
	URLHandler.RegisterHandler(MainPageHandler{}, "/")
	URLHandler.RegisterHandler(BugPageHandler{}, "/issues/")
	http.ListenAndServe(":8080", nil)
}

func getRootDir() string {
	dir := os.Getenv("PMIT")
	if dir != "" {
		return dir
	}

	wd, _ := os.Getwd()

	if dirinfo, err := os.Stat(wd + "/issues"); err == nil && dirinfo.IsDir() {
		return wd
	}

	// There's no environment variable and no issues
	// directory, so walk up the tree until we find one
	pieces := strings.Split(wd, "/")

	for i := len(pieces); i > 0; i -= 1 {
		dir := strings.Join(pieces[0:i], "/")
		if dirinfo, err := os.Stat(dir + "/issues"); err == nil && dirinfo.IsDir() {
			return dir
		}
	}
	return ""
}
