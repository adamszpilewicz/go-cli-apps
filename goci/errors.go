package main

import (
	"errors"
	"fmt"
)

var (
	ErrValidation = errors.New("validation failed")
)

type stepErr struct {
	step  string
	msg   string
	cause error
}

func (se *stepErr) Error() string {
	return fmt.Sprintf("step: %q: %s: cause: %v", se.step, se.msg, se.cause)
}

func (se *stepErr) Is(target error) bool {
	t, ok := target.(*stepErr)
	if !ok {
		return false
	}
	return se.step == t.step
}

func (se *stepErr) Unwrap() error {
	return se.cause
}
