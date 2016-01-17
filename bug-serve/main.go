package main

import (
	"github.com/driusan/GoWebapp/URLHandler"
	"net/http"
)

func main() {
	URLHandler.RegisterHandler(MainPageHandler{}, "/")
	URLHandler.RegisterHandler(SettingsHandler{}, "/settings")
	URLHandler.RegisterHandler(BugPageHandler{}, "/issues/")
	URLHandler.RegisterHandler(JSAssetHandler{}, "/js/")
	http.ListenAndServe(":8080", nil)
}

//go:generate babel --out-dir=js jsx
//go:generate go-bindata js/
