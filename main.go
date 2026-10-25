package main

import (
	"flag"
	"fmt"
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

//go:embed web/style.css
var styleCSS string

//go:embed imgs/a.jpg
var startIcon string

func init() {
	// Init flag
	flag.Usage = func() {
		fmt.Println(`leaflet-serve runs a server to expose a leaflet map with local GPX file
Usage:
	leaflet-serve path/to/a.gpx [other.gpx...]
`)
	}

	// Parse flag
	flag.Parse()
}

func main() {
	// Get positional args for GPX Files to serve
	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		fmt.Println("You must provide at least one gpx filepath to expose.")
		return
	}

	// Set up utils
	http.HandleFunc("/map.js", serveStatic(mapJs, http.Header{"Content-Type": []string{"text/javascript"}}))
	http.HandleFunc("/style.css", serveStatic(styleCSS, http.Header{"Content-Type": []string{"text/css"}}))
	// http.HandleFunc("/imgs/start.jpg", serveStatic(startIcon, http.Header{"Content-Type": []string{"image/jpeg"}}))

	// Set up main page
	http.HandleFunc("/", serveMainWithGPXFiles(leafletHTML, args))

	// Set up route to expose .gpx files
	for _, dir := range uniqueDirs(args) {
		route := filepath.Clean(ROUTE_PREFIX + dir + "{GPXFilePath}")
		http.HandleFunc(route, serveGPXFile(dir))
	}

	fmt.Println("Run web server on localhost:8000")
	http.ListenAndServe(":8000", nil)
}
