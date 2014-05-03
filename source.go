// This file defines the Source interface along with a way of merging
// multiple different sources into one Source.

package main

import (
	"errors"
	"fmt"
)

// A Source is capable of going to a given location and provides
// completion for a given prefix.
type Source interface {
	// Goto takes a destination string and returns a filepath to goto
	// or an error if this Source doesn't know how to get our file.
	Goto(dest string) (string, error)

	// Complete provides a list of strings representing completetions
	// this Source knows about that match the prefix or nil if there
	// are no completions.
	Complete(prefix string) []string
}

// mergedSources is a Source that merges multiple other Sources.
type mergedSource []Source

// Merge merges the given Sources into one. The merged sources Goto()
// function returns the first successful Goto of one of it's
// children. Complete() returns the union of all the child sources
// completions.
func Merge(s ...Source) Source {
	return mergedSource(s)
}

func (s mergedSource) Goto(str string) (string, error) {
	for _, source := range s {
		dir, err := source.Goto(str)
		if err == nil {
			return dir, err
		}
	}
	return "", errors.New(fmt.Sprintf("No such shortcut %s", str))
}

func (s mergedSource) Complete(prefix string) []string {
	var completions []string
	for _, source := range s {
		for _, completion := range source.Complete(prefix) {
			completions = append(completions, completion)
		}
	}
	return completions
}
