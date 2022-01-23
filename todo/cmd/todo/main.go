package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"todo"
)

var toDoName = ".todo_test.json"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool. \nDeveloped for The Pragmatic Bookshelf\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2020\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	}

	if os.Getenv("TODO_FILENAME") != "" {
		toDoName = os.Getenv("TODO_FILENAME")
	}

	// Parsing command line flags
	add := flag.Bool("add", false, " to be included in todo list")
	list := flag.Bool("list", false, "list all tasks")
	complete := flag.Int("complete", 0, "item to be completed")
	flag.Parse()

	l := &todo.List{}

	if err := l.Get(toDoName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *list:
		fmt.Println(l)
		//for _, item := range *l {
		//	if !item.Done {
		//		fmt.Println(item.Task)
		//	}
		//}
	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Save the new list
		if err := l.Save(toDoName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *add:
		// Add the
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		l.Add(t)

		// Save the new list
		if err := l.Save(toDoName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	default:
		// Invalid flag provided
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}
	s := bufio.NewScanner(r)
	s.Scan()
	if err := s.Err(); err != nil {
		return "", err
	}
	if len(s.Text()) == 0 {
		return "", fmt.Errorf("task cannot be blank")
	}
	return s.Text(), nil
}
