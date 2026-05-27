package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	_ "embed"
)

// ROUTE_PREFIX configures the route prefix to expose user files
const ROUTE_PREFIX = "/tmp/"

//go:embed web/leaflet.html
var leafletHtml string

//go:embed web/map.js
var mapJs string

func serveMapJs(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, mapJs)
	// http.ServeFile(w, r, "web/map.js")
}

func serveDir(dir string) func(http.ResponseWriter, *http.Request) {
	fmt.Println("serveDir: dir:", dir)
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(`serveDir: r.PathValue("gpxFilepath"):`, r.PathValue("gpxFilepath"))
		localFile := dir + r.PathValue("gpxFilepath")
		fmt.Println("Serve local file: ", localFile)
		http.ServeFile(w, r, localFile)
	}
}

func serveMainPage(localFilePaths []string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		funcs := template.FuncMap{}
		funcs["basename"] = filepath.Base
		t, err := template.New("").
			Funcs(funcs).
			Parse(leafletHtml)
		check(err)

		err = t.Execute(w, struct {
			GpxFilepath []string
		}{
			GpxFilepath: prefixSlice(localFilePaths, ROUTE_PREFIX)})
		check(err)
	}
}

func main() {
	// Set up main page, template filled
	http.HandleFunc("/", serveMainPage(os.Args[1:]))

	// Set up map.js
	http.HandleFunc("/map.js", serveMapJs)

	// Set up route to expose .gpx files
	dirs := dirsToExpose(os.Args[1:])
	for _, dir := range dirs {

		// TODO: fix: ./ -> "" not working
		// it redirs to /tmp/tmp/tmp...
		routeSuffix := filepath.Clean(dir + "{gpxFilepath}")
		dirToServe, _ := filepath.Split(routeSuffix)
		// dirToServe := dir
		// routeSuffix := dir + "{gpxFilepath}"
		route := ROUTE_PREFIX + routeSuffix
		fmt.Println("Create route to: ", route)

		fmt.Println("main: dir:", dir)
		fmt.Println("main: dirToServe:", dirToServe)
		http.HandleFunc(route, serveDir(dirToServe))
	}

	fmt.Println("Run web server on localhost:8000")
	http.ListenAndServe(":8000", nil)
}
