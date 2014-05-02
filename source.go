package main

type source interface {
	Goto(str string) (string, error)
	Complete(string) []string
}
