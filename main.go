package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

const whitespace = " \t"

func main() {
	log.SetFlags(log.Lshortfile)

	var err error
	if len(os.Args) < 2 { // Go home
		fmt.Println(os.Getenv("HOME"))
	} else if len(os.Args) == 2 && os.Args[1] != "-complete" { // Go somewhere
		err = gocd()
	} else if len(os.Args) <= 3 && os.Args[1] == "-complete" { // Go complete
		err = complete()
	} else {
		err = errors.New("Improper usage...")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func gocd() error {
	source := getSources()
	dir, err := source.Goto(os.Args[1])
	if err != nil {
		return err
	}
	fmt.Println(dir)
	return nil
}

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

func getSources() source {
	return NewFileSource(os.Getenv("HOME") + "/.goto")
}
