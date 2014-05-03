// This file defines a Source that goes to and completes paths for Go.

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// NewGoSource returns a Source that goes to and completes paths for
// Go. It respects the current $GOROOT and the directories on the
// $GOPATH just as `go build` does.
func NewGoSource() Source {
	sources := mergedSource{NewDirSource(os.Getenv("GOROOT") +
		string(os.PathSeparator) + "src" + string(os.PathSeparator) + "pkg")}

	for _, source := range filepath.SplitList(os.Getenv("GOPATH")) {
		sources = append(sources, NewDirSource(source+string(os.PathSeparator)+"src"))
	}

	return sources
}

// NewDirSource returns a Source that provides recursive completions
// of subdirectories of the given directory.
func NewDirSource(dir string) Source {
	return &dirSource{dir}
}

type dirSource struct {
	dir string
}

// Goto goes to the given subdirectory if it exists on the filesystem.
func (s *dirSource) Goto(dest string) (string, error) {
	dir := s.dir + string(os.PathSeparator) + dest
	_, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", errors.New(fmt.Sprintf("No such shortcut %s", dest))
	}
	return dir, nil
}

// Complete returns recursive completion of all (sub)+directories of
// the current directory source.
func (s *dirSource) Complete(prefix string) []string {
	n := strings.LastIndex(prefix, string(os.PathSeparator))
	if n == -1 {
		return completionsForDir(s.dir, prefix)
	} else {
		dir := s.dir + string(os.PathSeparator) + prefix[:n]
		var completions []string
		for _, completion := range completionsForDir(dir, prefix[n+1:]) {
			completions = append(completions, prefix[:n]+string(os.PathSeparator)+completion)
		}
		return completions
	}
}

func completionsForDir(dir, prefix string) []string {
	paths, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil
	}

	var completions []string
	for _, path := range paths {
		if path.IsDir() && strings.HasPrefix(path.Name(), prefix) {
			completions = append(completions, path.Name()+string(os.PathSeparator))
		}
	}
	return completions
}
