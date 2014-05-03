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
	val, err := s.expand(str)
	if err != nil {
		return "", errors.New(fmt.Sprintf("No such shortcut %s", str))
	}
	return val, err
}

func (s *FileSource) Complete(prefix string) []string {
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

func (s *FileSource) expand(str string) (string, error) {
	for key := range s.shortcuts {
		if strings.HasPrefix(str, key) {
			return s.shortcuts[key] + string(os.PathSeparator) + strings.TrimPrefix(str, key), nil
		}
	}
	return "", errors.New(fmt.Sprintf("can't", str))
}
