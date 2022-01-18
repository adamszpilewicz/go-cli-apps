package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

type executer interface {
	execute() (string, error)
}

func main() {
	proj := flag.String("proj", ".", "project directory")
	flag.Parse()

	if err := run(*proj, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(proj string, out io.Writer) error {
	if proj == "" {
		return fmt.Errorf("direcotry has not been provided: %w", ErrValidation)
	}

	pipeline := make([]executer, 4)
	pipeline[0] = newStep(
		"go build",
		"go",
		"Go Build: SUCCESS",
		proj,
		[]string{"build", ".", "errors"},
	)
	pipeline[1] = newStep(
		"go test",
		"go",
		"Go Test: SUCCESS",
		proj,
		[]string{"test", "-v"},
	)
	pipeline[2] = newExceptionStep(
		"go fmt",
		"gofmt",
		"Gofmt: SUCCESS",
		proj,
		[]string{"-l", "."},
	)
	pipeline[3] = newTimeoutStep("git push",
		"git",
		"Git Push: SUCCESS",
		proj,
		[]string{"push", "origin", "master"}, 10*time.Second,
	)

	for _, s := range pipeline {
		msg, err := s.execute()
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(out, msg)
		if err != nil {
			return err
		}
	}
	return nil
	//args := []string{"build", ".", "errors"}
	//cmd := exec.Command("go", args...)
	//cmd.Dir = proj
	//if err := cmd.Run(); err != nil {
	//	return &stepErr{step: "go build", msg: "go build failed", cause: err}
	//}
	//_, err := fmt.Fprintln(out, "Go Build: SUCCESS")
	//return err
}
