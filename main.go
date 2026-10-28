package main

import (
	"embed"
	"flag"
	"fmt"
	"net/http"
	"path/filepath"

	_ "embed"
)

// ROUTE_PREFIX configures the route prefix to expose user files
const ROUTE_PREFIX = "/gpx/"

//go:embed web/leaflet.html
var leafletHTML string

//go:embed web/map.js
var mapJs string

//go:embed web/style.css
var styleCSS string

//go:embed imgs/*
var ficons embed.FS

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
	http.HandleFunc("/map.js", serveString(mapJs, http.Header{"Content-Type": []string{"text/javascript"}}))
	http.HandleFunc("/style.css", serveString(styleCSS, http.Header{"Content-Type": []string{"text/css"}}))

	http.HandleFunc("/imgs/flag-start.png",
		serveFile("imgs/flag-start.png", http.Header{"Content-Type": []string{"image/png"}}))
	http.HandleFunc("/imgs/flag-end.png",
		serveFile("imgs/flag-end.png", http.Header{"Content-Type": []string{"image/png"}}))

	http.HandleFunc("/imgs/person-shelter.svg",
		serveFile("imgs/person-shelter.svg", http.Header{"Content-Type": []string{"image/svg+xml"}}))

	http.HandleFunc("/imgs/building.svg",
		serveFile("imgs/building.svg", http.Header{"Content-Type": []string{"image/svg+xml"}}))
	http.HandleFunc("/imgs/water.svg",
		serveFile("imgs/water.svg", http.Header{"Content-Type": []string{"image/svg+xml"}}))
	http.HandleFunc("/imgs/mountain-sun.svg",
		serveFile("imgs/mountain-sun.svg", http.Header{"Content-Type": []string{"image/svg+xml"}}))
	http.HandleFunc("/imgs/eye.svg",
		serveFile("imgs/eye.svg", http.Header{"Content-Type": []string{"image/svg+xml"}}))

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
