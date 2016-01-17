package main

import (
	"github.com/driusan/GoWebapp/URLHandler"
	"net/http"
)

type JSAssetHandler struct {
	URLHandler.DefaultHandler
}

func (s JSAssetHandler) Get(r *http.Request, p map[string]interface{}) (string, error) {
	data, err := Asset(r.URL.Path[1:])
	if err == nil {
		return string(data), nil

	}
	return "File Not found", URLHandler.NotFoundError{}
}
