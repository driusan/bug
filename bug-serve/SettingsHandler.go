package main

import (
	"encoding/json"
	"github.com/driusan/GoWebapp/URLHandler"
	"github.com/driusan/bug/bugs"
	"net/http"
)

type SettingsHandler struct {
	URLHandler.DefaultHandler
}

func (s SettingsHandler) Get(r *http.Request, p map[string]interface{}) (string, error) {
	settings := struct {
		Title     string
		Directory string
	}{bugs.GetRootDir().GetShortName().ToTitle(), string(bugs.GetRootDir())}
	retVal, _ := json.Marshal(settings)
	return string(retVal), nil
}
