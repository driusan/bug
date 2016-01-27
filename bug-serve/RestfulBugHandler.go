package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/driusan/GoWebapp/URLHandler"
	"github.com/driusan/bug/bugs"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

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
func parseURL(u *url.URL) []string {
	if u.Path == "/issues" || u.Path == "/issues/" {
		return []string{}
	}
	bugURL := strings.TrimPrefix(u.Path, "/issues/")
	bugURL = strings.TrimRight(bugURL, "/")

	return strings.Split(bugURL, "/")

}

func (m BugPageHandler) Get(r *http.Request, extras map[string]interface{}) (string, error) {
	switch urlChunks := parseURL(r.URL); len(urlChunks) {
	case 0:
		return getBugList()
	case 1:
		b, err := bugs.LoadBugByDirectory(urlChunks[0])
		if err != nil {
			return "", URLHandler.NotFoundError{}
		}
		return b.ToJSONString()
	default:
		// Must be > 1, so check what field we're editing..
		// Need to load the bug since this is a different
		// case statement..
		b, err := bugs.LoadBugByDirectory(urlChunks[0])
		if err != nil {
			return "", URLHandler.NotFoundError{}
		}
		switch urlChunks[1] {
		case "Description":
			return b.Description(), nil
		default:
			// No known field, so return a 404 error.
			return "", URLHandler.NotFoundError{}
		}
	}
	return "", URLHandler.NotFoundError{}
}

func (m BugPageHandler) Put(r *http.Request, extras map[string]interface{}) (string, error) {
	urlChunks := parseURL(r.URL)
	switch len(urlChunks) {
	case 0:
		return "", URLHandler.BadRequestError{}
	case 1:
		// This should eventually be supported if a JSON string
		// is PUT, but for now we only support PUT to a specific
		// field
		return "", URLHandler.BadRequestError{}
	}

	b, err := bugs.LoadBugByDirectory(urlChunks[0])
	if err != nil {
		// This should create the issue instead, but for now
		// only updating fields is supported
		return "", URLHandler.BadRequestError{}
	}
	switch urlChunks[1] {
	case "Description":
		if val, err := ioutil.ReadAll(r.Body); err == nil {
			b.SetDescription(string(val))
			return "", nil
		} else {
			panic(err.Error())
		}

	}
	return "", URLHandler.BadRequestError{}
}

func (m BugPageHandler) Delete(r *http.Request, extras map[string]interface{}) (string, error) {
	urlChunks := parseURL(r.URL)
	if len(urlChunks) != 1 {
		return "", URLHandler.BadRequestError{}
	}

	b, err := bugs.LoadBugByDirectory(urlChunks[0])
	if err != nil {
		return "", URLHandler.NotFoundError{}
	}

	if err := b.Remove(); err != nil {
		panic("Could not delete bug.")
	}
	return "", nil
}
func (m BugPageHandler) ETag(u *url.URL, o map[string]interface{}) (URLHandler.ETag, error) {
	urlChunks := parseURL(u)
	fmt.Printf("Calculating ETag for %s => %s\n", u, urlChunks)
	if len(urlChunks) == 0 {
		return URLHandler.ETag(""), nil
	}

	b, err := bugs.LoadBugByDirectory(urlChunks[0])
	if err != nil {
		// If the bug doesn't exist, it's not an error, but there
		// shouldn't be an ETag either..
		return "", nil
	}

	h := sha1.New()
	io.WriteString(h, b.Title(""))
	io.WriteString(h, "--")
	io.WriteString(h, b.Description())
	io.WriteString(h, "--")
	io.WriteString(h, b.Identifier())
	io.WriteString(h, "--")
	io.WriteString(h, b.Status())
	io.WriteString(h, "--")
	io.WriteString(h, b.Milestone())
	io.WriteString(h, "--")
	io.WriteString(h, b.Priority())
	io.WriteString(h, "--")
	io.WriteString(h, b.GetDirectory().LastModified().Format(time.UnixDate))
	return URLHandler.ETag(fmt.Sprintf("%x", h.Sum(nil))), nil
}

func (m BugPageHandler) LastModified(u *url.URL, o map[string]interface{}) time.Time {
	urlChunks := parseURL(u)
	if len(urlChunks) == 0 {
		return URLHandler.UnknownMTime
	}

	b, err := bugs.LoadBugByDirectory(urlChunks[0])
	if err != nil {
		return URLHandler.UnknownMTime
	}
	return b.GetDirectory().LastModified()

}
