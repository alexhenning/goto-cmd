package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func NewGoSource() source {
	sources := sources{NewDirSource(os.Getenv("GOROOT") +
		string(os.PathSeparator) + "src" + string(os.PathSeparator) + "pkg")}

	for _, source := range filepath.SplitList(os.Getenv("GOPATH")) {
		sources = append(sources, NewDirSource(source+string(os.PathSeparator)+"src"))
	}

	return sources
}

type DirSource struct {
	dir string
}

func NewDirSource(dir string) source {
	return &DirSource{dir}
}

func (s *DirSource) Goto(str string) (string, error) {
	d := s.dir + string(os.PathSeparator) + str
	_, err := ioutil.ReadDir(d)
	if err != nil {
		return "", errors.New(fmt.Sprintf("No such shortcut %s", str))
	}
	return d, nil
}

func (s *DirSource) Complete(prefix string) []string {
	// Get directory
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

	// Completions for directory
	var completions []string
	for _, path := range paths {
		if path.IsDir() && strings.HasPrefix(path.Name(), prefix) {
			completions = append(completions, path.Name()+string(os.PathSeparator))
		}
	}
	return completions
}
