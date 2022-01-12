package main_test

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	binName  = "todo"
	fileName = ".todo_test.json"
)

// TestMain executing extra setup before your
// tests is by using the TestMain() function.
// This function helps you control the extra tasks
// required to set up or tear down the resources
// necessary for testing, keeping your test cases
// tidy and consistent.
func TestMain(m *testing.M) {
	log.Println("---> creating file")
	_, err := os.Create(".todo_test.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot create file, err: %v", err)
		os.Exit(1)
	}

	log.Println("---> building executable...")
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	build := exec.Command("go", "build", "-o", binName)
	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "cannot build executable %s: error - %s", binName, err)
		os.Exit(1)
	}

	log.Println("---> running tests...")
	result := m.Run()

	log.Println("---> cleaning up...")
	os.Remove(binName)
	os.Remove(fileName)

	os.Exit(result)
}

func TestToDoCLI(t *testing.T) {
	task := "test task number 1"
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cmdPath := filepath.Join(dir, binName)

	t.Run("AddNewTask", func(t *testing.T) {
		//cmd := exec.Command(cmdPath, strings.Split(task, " ")...)
		cmd := exec.Command(cmdPath, "-task", task)

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ListTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		expected := fmt.Sprint("[ ] 1: ", task, "\n\n")
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})
}
