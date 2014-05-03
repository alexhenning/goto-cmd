package main

import (
	"errors"
	"fmt"
)

type source interface {
	Goto(str string) (string, error)
	Complete(string) []string
}

type sources []source

func Merge(s ...source) source {
	return sources(s)
}

func (s sources) Goto(str string) (string, error) {
	for _, source := range s {
		dir, err := source.Goto(str)
		if err == nil {
			return dir, err
		}
	}
	return "", errors.New(fmt.Sprintf("No such shortcut %s", str))
}

func (s sources) Complete(prefix string) []string {
	var completions []string
	for _, source := range s {
		for _, completion := range source.Complete(prefix) {
			completions = append(completions, completion)
		}
	}
	return completions
}
