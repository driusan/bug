package main

import (
	"fmt"
	"github.com/driusan/GoWebapp/HTMLPageRenderer"
	"github.com/driusan/GoWebapp/URLHandler"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"net/http"
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
	bugDir := bugs.GetRootDir() + r.URL.Path
	b := bugs.Bug{}
	b.LoadBug(bugs.Directory(bugDir))

	page := BugRenderer{Bug: b}
	page.RootElement = "RBugPage"
	page.Title = b.Title
	page.JSFiles = []string{
	// Bootstrap JS
	//"https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/js/bootstrap.min.js",
	// React JS
	"https://cdnjs.cloudflare.com/ajax/libs/react/0.14.3/react.js",
	"https://cdnjs.cloudflare.com/ajax/libs/react/0.14.3/react-dom.js",
	"/js/BugPage.js",
}
	page.CSSFiles = []string{
"https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css",
	"https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap-theme.min.css"}
	return HTMLPageRenderer.Render(page), nil

}

func main() {
	URLHandler.RegisterHandler(MainPageHandler{}, "/")
	URLHandler.RegisterHandler(BugPageHandler{}, "/issues/")
	URLHandler.RegisterStaticHandler("/js/", "./js")
	http.ListenAndServe(":8080", nil)
}
