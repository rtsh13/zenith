package server

import "fmt"

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
