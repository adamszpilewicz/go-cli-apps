package main

import (
	"fmt"
	"os"
	"strings"
	"todo"
)

const toDoName = ".todo_test.json"

func main() {
	l := &todo.List{}

	if err := l.Get(toDoName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case len(os.Args) == 1:
		for _, item := range *l {
			fmt.Println(item.Task)
		}
	default:
		item := strings.Join(os.Args[1:], " ")
		l.Add(item)

		if err := l.Save(toDoName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
