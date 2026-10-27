package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
)

/*
serveString serves static files as *.js script, *.css stylesheet, ...
*/
func serveString(content string, header http.Header) func(http.ResponseWriter, *http.Request) {
	var h map[string][]string = header // convert header to fetch keys
	return func(w http.ResponseWriter, r *http.Request) {
		// Set response header
		for key := range h {
			w.Header().Set(key, header.Get(key))
		}

		// Write body
		fmt.Fprintf(w, content)
	}
}
func serveFile(path string, header http.Header) func(http.ResponseWriter, *http.Request) {
	var h map[string][]string = header // convert header to fetch keys
	return func(w http.ResponseWriter, r *http.Request) {
		// Set response header
		for key := range h {
			w.Header().Set(key, header.Get(key))
		}

		// Write body
		http.ServeFile(w, r, path)
	}
}

/*
serverGPCFile serves GPX files on "pattern/{GPXFilePath}"
*/
func serveGPXFile(pattern string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
