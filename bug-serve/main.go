package main

import (
	"encoding/json"
	//	"fmt"
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

type SettingsHandler struct {
	URLHandler.DefaultHandler
}
type BugListRenderer struct {
	HTMLPageRenderer.ReactPage
}

type BugRenderer struct {
	HTMLPageRenderer.ReactPage
	Bug bugs.Bug
}

func (s SettingsHandler) Get(r *http.Request, p map[string]interface{}) (string, error) {
	settings := struct {
		Directory string
	}{string(bugs.GetRootDir())}
	retVal, _ := json.Marshal(settings)
	return string(retVal), nil
}
func (m MainPageHandler) Get(r *http.Request, p map[string]interface{}) (string, error) {
	page := BugListRenderer{}
	page.Title = "Open Issues"
	page.JSFiles = []string{
		// Bootstrap
		//"https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/js/bootstrap.min.js",
		// React
		"https://cdnjs.cloudflare.com/ajax/libs/react/0.14.3/react.js",
		"https://cdnjs.cloudflare.com/ajax/libs/react/0.14.3/react-dom.js",
		"/js/BugApp.js",
		"/js/BugList.js",
		"/js/BugPage.js",
	}
	page.CSSFiles = []string{
		// Bootstrap
		"https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css",
		"https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap-theme.min.css",
	}
	page.RootElement = "BugApp"

	return HTMLPageRenderer.Render(page), nil
}

func getBugList() (string, error) {
	issues, _ := ioutil.ReadDir(string(bugs.GetRootDir()) + "/issues")

	var issuesSlice []string

	for _, issue := range issues {
		issuesSlice = append(issuesSlice, issue.Name())
	}

	retVal, _ := json.Marshal(issuesSlice)
	return string(retVal), nil
}
func (m BugPageHandler) Get(r *http.Request, extras map[string]interface{}) (string, error) {
	if r.URL.Path == "/issues" || r.URL.Path == "/issues/" {
		return getBugList()
	}
	bugDir := string(bugs.GetRootDir()) + r.URL.Path
	b := bugs.Bug{}
	b.LoadBug(bugs.Directory(bugDir))

	switch r.URL.Query().Get("format") {
	case "json":
		bJSON, _ := json.Marshal(b)
		return string(bJSON), nil
	default:
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

}

func main() {
	URLHandler.RegisterHandler(MainPageHandler{}, "/")
	URLHandler.RegisterHandler(SettingsHandler{}, "/settings")
	URLHandler.RegisterHandler(BugPageHandler{}, "/issues/")
	URLHandler.RegisterStaticHandler("/js/", "./js")
	http.ListenAndServe(":8080", nil)
}
