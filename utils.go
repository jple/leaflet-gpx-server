package main

import (
	"fmt"
	"path/filepath"
	"slices"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

/*
uniqueDirs returns unique slice of directories to expose
example:

	input : /tmp/A.gpx /tmp/B.gpx /tmp/dir/C.gpx
	output : /tmp/ /tmp/dir/ (prevent dupplicates /tmp/ dir)
*/
func uniqueDirs(paths []string) []string {
	var dirs []string
	for _, p := range paths {
		d := filepath.Dir(p) + "/"
		if !slices.Contains(dirs, d) {
			dirs = append(dirs, d)
		}
	}
	return dirs
}

func prefixSlice(s []string, prefix string) []string {
	ss := make([]string, len(s)) // copy to prevent modifying original content
	for i := range ss {
		ss[i] = prefix + s[i]
	}
	return ss
}
