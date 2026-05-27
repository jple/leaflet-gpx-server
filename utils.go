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
dirsToExpose is a helper function returning a unique slice of directories to expose
It is used to parse user input
*/
func dirsToExpose(paths []string) []string {
	var dirs []string
	for _, p := range paths {
		d := filepath.Dir(p) + "/"
		if !slices.Contains(dirs, d) {
			dirs = append(dirs, d)
		}
	}
	return dirs
}

func prefixSlice(sl []string, prefix string) []string {
	for i := range sl {
		sl[i] = prefix + sl[i]
	}
	return sl
}
