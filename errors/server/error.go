package server

import (
	"fmt"
	"strings"
)

type ProtocolError struct {
	Message string
}

func (p ProtocolError) Error() string {
	return fmt.Sprintf("-ERR Protocol error: %s", p.Message)
}

type ScannerError struct {
	Message string
}

func (s ScannerError) Error() string {
	return fmt.Sprintf("-ERR Input Scan error: %s", s.Message)
}

type UnknownCommand struct {
	Command string
	Args    []string
}

func (u UnknownCommand) Error() string {
	return fmt.Sprintf("(error) ERR unknown command '%v', with args beginning with: %v", u.Command, u.Args)
}

type MultipleErrors struct {
	Errors []error
}

func (m MultipleErrors) Error() string {
	errors := []string{}

	for _, e := range m.Errors {
		errors = append(errors, e.Error())
	}

	return strings.Join(errors, ";")
}
