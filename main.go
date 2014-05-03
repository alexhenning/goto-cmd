package main

import (
	"fmt"
	"os"
)

// main provides either a directory to cd into, a list of completion
// or usage instructions.
func main() {
	var err error
	if len(os.Args) < 2 {
		fmt.Println(os.Getenv("HOME"))
	} else if os.Args[1] == "-h" || os.Args[1] == "-help" || os.Args[1] == "--help" {
		fmt.Fprintln(os.Stderr, Usage)
	} else if len(os.Args) == 2 && os.Args[1] != "-complete" {
		err = gocd()
	} else if len(os.Args) <= 3 && os.Args[1] == "-complete" {
		err = complete()
	} else {
		fmt.Fprintln(os.Stderr, Usage)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

// gocd prints the directory to cd to stdout.
func gocd() error {
	source := getSources()
	dir, err := source.Goto(os.Args[1])
	if err != nil {
		return err
	}
	fmt.Println(dir)
	return nil
}

// complete print a list of completions to stdout, one per line.
func complete() error {
	source := getSources()

	var prefix string
	if len(os.Args) == 3 {
		prefix = os.Args[2]
	}

	for _, completion := range source.Complete(prefix) {
		fmt.Println(completion)
	}

	return nil
}

// getSources returns the sources for completing from the users .goto
// and the completions for the current Go source tree.
func getSources() Source {
	return Merge(NewFileSource(os.Getenv("HOME")+"/.goto"), NewGoSource())
}
