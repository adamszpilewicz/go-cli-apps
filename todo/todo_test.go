package todo_test

import (
	"log"
	"testing"
	"todo"
)

func TestList_Add(t *testing.T) {
	l := todo.List{}
	taskName := "added task"
	l.Add(taskName)
	if l[0].Task != taskName {
		t.Errorf("expected %q, got %q", taskName, l[0].Task)
	}
}

func TestList_Complete(t *testing.T) {
	l := todo.List{}
	taskName := "completed task"
	l.Add(taskName)
	if l[0].Done {
		t.Errorf("newly created task should not have status completed: %v", l[0].Done)
	}

	err := l.Complete(1)
	if err != nil {
		log.Println(err)
		return
	}
	if !l[0].Done {
		t.Errorf("expected done=true, got %v", l[0].Done)
	}
}
