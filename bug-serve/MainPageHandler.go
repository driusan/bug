package main

import (
	"github.com/driusan/GoWebapp/HTMLPageRenderer"
	"github.com/driusan/GoWebapp/URLHandler"
	"net/http"
)

type MainPageHandler struct {
	URLHandler.DefaultHandler
}

func (m MainPageHandler) Get(r *http.Request, p map[string]interface{}) (string, error) {
	page := HTMLPageRenderer.ReactPage{
		Title: "Open Issues",
		JSFiles: []string{
			// Bootstrap
			//"https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/js/bootstrap.min.js",
			// React
			"https://cdnjs.cloudflare.com/ajax/libs/react/0.14.3/react.js",
			"https://cdnjs.cloudflare.com/ajax/libs/react/0.14.3/react-dom.js",
			"/js/BugApp.js",
			"/js/BugList.js",
			"/js/BugPage.js",
		},
		CSSFiles: []string{
			// Bootstrap
			"https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css",
			"https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap-theme.min.css",
		},
		RootElement: "BugApp",
	}
	return HTMLPageRenderer.Render(page), nil
}
