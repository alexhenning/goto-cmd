// goto [-complete] dir
//
// Goes to the given directory if possible otherwise. If "-complete" is
// present, instead of returning the directory to cd into, it returns a
// list of the possible completions.
//
// The format for the ~/.goto file is a pair of labels and
// corresponding paths. Blank lines and lines beginning with `#` are
// ignored.
//
// ~/.goto:
// # label path
// ~go     ~/Dropbox/Programming/go
// me      ~/Dropbox/Programming/go/src/github.com/alexhenning
package main

var Usage = `goto [-complete] dir

Goes to the given directory if possible otherwise. If "-complete" is
present, instead of returning the directory to cd into, it returns a
list of the possible completions.`
