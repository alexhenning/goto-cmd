package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const whitespace = " \t"

func main() {
	log.SetFlags(log.Lshortfile)
	if len(os.Args) < 2 {
		fmt.Println(os.Getenv("HOME"))
		return
	}

	shortcuts := loadShortcuts(os.Getenv("HOME") + "/.goto")
	if dir, ok := shortcuts[os.Args[1]]; ok {
		fmt.Println(dir)
	} else {
		log.Fatalf("No such shortcut %s", os.Args[1])
	}
}

func loadShortcuts(file string) map[string]string {
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
	return shortcuts
}
