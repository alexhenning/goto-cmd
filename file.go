package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type FileSource struct {
	shortcuts map[string]string
}

func NewFileSource(file string) source {
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
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
	return &FileSource{shortcuts}
}

func (s *FileSource) Goto(str string) (string, error) {
	dir, ok := s.shortcuts[str]
	if !ok {
		return "", errors.New(fmt.Sprintf("No such shortcut %s", str))
	}
	return dir, nil
}

func (s *FileSource) Complete(prefix string) []string {
	var completions []string
	for key := range s.shortcuts {
		if strings.HasPrefix(key, prefix) {
			completions = append(completions, key)
		}
	}
	return completions
}
