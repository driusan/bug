package main

import (
	"fmt"
	"github.com/driusan/GoWebapp/HTMLPageRenderer"
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

type BugListRenderer struct {
	HTMLPageRenderer.Page
}

func (b BugListRenderer) GetBody() string {
	issues, _ := ioutil.ReadDir(bugs.GetRootDir() + "/issues")

	ret := "<h2>" + b.Title + "</h2><ol>"
	for _, issue := range issues {
		var dir bugs.Directory = bugs.Directory(issue.Name())
		ret += fmt.Sprintf("<li><a href=\"/issues/%s\">%s</a></li>\n", (dir), dir.ToTitle())
	}
	ret += "</ol>"

	return ret
}

type BugRenderer struct {
	HTMLPageRenderer.ReactPage
	Bug bugs.Bug
}
func (m MainPageHandler) Get(*http.Request, map[string]interface{}) (string, error) {

	page := BugListRenderer{}
	page.Title = "Open Issues"
	page.JSFiles = []string{"https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/js/bootstrap.min.js"}
	page.CSSFiles = []string{
"https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css",
	"https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap-theme.min.css"}

	return HTMLPageRenderer.Render(page), nil
}
func (m BugPageHandler) Get(r *http.Request, extras map[string]interface{}) (string, error) {
	bugDir := getRootDir() + r.URL.Path
	b := bugs.Bug{}
	b.LoadBug(bugs.Directory(bugDir))

	page := BugRenderer{Bug: b}
	page.RootElement = "abc"
	page.Title = b.Title
	page.JSFiles = []string{"https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/js/bootstrap.min.js"}
	page.CSSFiles = []string{
"https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css",
	"https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap-theme.min.css"}
	return HTMLPageRenderer.Render(page), nil

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
