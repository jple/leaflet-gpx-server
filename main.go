package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	_ "embed"
)

// ROUTE_PREFIX configures the route prefix to expose user files
const ROUTE_PREFIX = "/tmp/"

//go:embed web/leaflet.html
var leafletHTML string

//go:embed web/map.js
var mapJs string

func serveMapJs(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, mapJs)
	// http.ServeFile(w, r, "web/map.js")
}

func servePattern(pattern string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println("request:", r.Pattern) // not working
		http.ServeFile(w, r, pattern+r.PathValue("GPXFilePath"))
	}
}

func writeTemplateWithGPXFiles(w io.Writer, tmpl string, paths []string) {
	funcs := template.FuncMap{}
	funcs["basename"] = filepath.Base

	t, err := template.New("").
		Funcs(funcs).
		Parse(tmpl)
	check(err)

	err = t.Execute(w, struct {
		GPXFilepath []string
	}{
		GPXFilepath: paths})
	check(err)
}

func serveMainWithGPXFiles(HTMLmain string, GPXFilepaths []string) func(w http.ResponseWriter, r *http.Request) {
	// format GPX filepaths
	GPXFilesToServe := prefixSlice(GPXFilepaths, ROUTE_PREFIX)
	for i := range GPXFilepaths {
		GPXFilesToServe[i] = filepath.Clean(GPXFilesToServe[i])
	}

	// Create template and write response with it
	return func(w http.ResponseWriter, r *http.Request) {
		writeTemplateWithGPXFiles(w, HTMLmain, GPXFilesToServe)
	}
}

func main() {
	// Parse positional args for GPX Files to serve
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("You must provide at least one gpx filepath to expose")
		return
	}

	// Set up utils
	http.HandleFunc("/map.js", serveMapJs)

	// Set up main page
	http.HandleFunc("/", serveMainWithGPXFiles(leafletHTML, args))

	// Set up route to expose .gpx files
	for _, dir := range uniqueDirs(args) {
		route := filepath.Clean(ROUTE_PREFIX + dir + "{GPXFilePath}")
		http.HandleFunc(route, servePattern(dir))
	}

	fmt.Println("Run web server on localhost:8000")
	http.ListenAndServe(":8000", nil)
}
