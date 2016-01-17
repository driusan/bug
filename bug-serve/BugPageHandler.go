package main

import (
	"encoding/json"
	"github.com/driusan/GoWebapp/HTMLPageRenderer"
	"github.com/driusan/GoWebapp/URLHandler"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"net/http"
	"strings"
)

type BugRenderer struct {
	HTMLPageRenderer.ReactPage
	Bug bugs.Bug
}
type BugPageHandler struct {
	URLHandler.DefaultHandler
}

func getBugList() (string, error) {
	issues, _ := ioutil.ReadDir(string(bugs.GetRootDir()) + "/issues")

	var issuesSlice []string

	for _, issue := range issues {
		if issue.IsDir() {
			issuesSlice = append(issuesSlice, issue.Name())
		}
	}

	retVal, _ := json.Marshal(issuesSlice)
	return string(retVal), nil
}

func (m BugPageHandler) Get(r *http.Request, extras map[string]interface{}) (string, error) {
	if r.URL.Path == "/issues" || r.URL.Path == "/issues/" {
		return getBugList()
	}
	// Strip off the "/issues/"
	bugURL := r.URL.Path[8:]

	// See if we're getting the whole bug JSON or
	// a field from it.
	pieces := strings.Split(bugURL, "/")

	b, err := bugs.LoadBugByDirectory(pieces[0])
	if err != nil {
		return "", URLHandler.NotFoundError{}
	}

	bJSONStruct := struct {
		Identifier  string `json:",omitempty"`
		Title       string
		Description string
		Status      string   `json:",omitempty"`
		Priority    string   `json:",omitempty"`
		Milestone   string   `json:",omitempty"`
		Tags        []string `json:",omitempty"`
	}{
		Identifier:  b.Identifier(),
		Title:       b.Title(""),
		Description: b.Description(),
		Status:      b.Status(),
		Priority:    b.Priority(),
		Milestone:   b.Milestone(),
		Tags:        b.StringTags(),
	}

	bJSON, _ := json.Marshal(bJSONStruct)
	return string(bJSON), nil
	/*
		switch r.URL.Query().Get("format") {
		case "json":
		default:
			page := BugRenderer{Bug: b}
			page.RootElement = "RBugPage"
			page.Title = b.Title("")
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
	*/

}
