// This file defines a Source that goes to and completes paths
// provided by a file.

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// NewFileSource returns a Source that goes to and completes path
// given by the file. The file format is given in the `doc.go` file.
func NewFileSource(file string) Source {
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return &fileSource{}
	}

	shortcuts := make(map[string]string)
	for i, line := range strings.Split(string(contents), "\n") {
		if strings.HasPrefix(line, "#") || strings.Trim(line, whitespace) == "" {
			continue
		}
		fields := strings.Fields(strings.Trim(line, whitespace))
		if len(fields) != 2 {
			log.Fatalf(".goto has more than two fields on line %d", i)
		}
		if strings.HasPrefix(fields[1], "~"+string(os.PathSeparator)) {
			fields[1] = os.Getenv("HOME") + strings.TrimPrefix(fields[1], "~")
		}
		shortcuts[fields[0]] = fields[1]
	}
	return &fileSource{shortcuts}
}

type fileSource struct {
	shortcuts map[string]string
}

// Goto goes to the expanded destination if it exists on the filesystem.
func (s *fileSource) Goto(dest string) (string, error) {
	dir, err := s.expand(dest)
	if err != nil {
		return "", errors.New(fmt.Sprintf("No such shortcut %s", dest))
	}
	return dir, nil
}

// Complete returns recursive completion of all (sub)+directories of
// the expanded prefix.
func (s *fileSource) Complete(prefix string) []string {
	var completions []string
	for key := range s.shortcuts {
		if strings.HasPrefix(key, prefix) {
			completions = append(completions, key+string(os.PathSeparator))
		} else if strings.HasPrefix(prefix, key) {
			val, err := s.expand(prefix)
			n := strings.LastIndex(val, string(os.PathSeparator))
			m := strings.LastIndex(prefix, string(os.PathSeparator))
			if err != nil || n == -1 || m == -1 {
				continue
			}
			for _, completion := range completionsForDir(val[:n], val[n+1:]) {
				completions = append(completions, prefix[:m+1]+completion)
			}
		}
	}
	return completions
}

func (s *fileSource) expand(str string) (string, error) {
	path, length := "", -1
	for key := range s.shortcuts {
		if strings.HasPrefix(str, key) && len(key) > length {
			path, length = s.shortcuts[key]+strings.TrimPrefix(str, key), len(key)
		}
	}
	if length == -1 {
		return "", errors.New(fmt.Sprintf("can't", str))
	}
	return path, nil
}

const whitespace = " \t"
