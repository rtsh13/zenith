package zenith

import "fmt"

type InvalidArgs struct {
	Command string
}

func (i InvalidArgs) Error() string {
	return fmt.Sprintf("(error) ERR wrong number of arguments for '%v' command", i.Command)
}

type UnknownCommand struct {
	Command string
	Args    []string
}

func (u UnknownCommand) Error() string {
	return fmt.Sprintf("(error) ERR unknown command '%v', with args beginning with: %v", u.Command, u.Args)
}

type ProtocolError struct {
	Message string
}

func (p ProtocolError) Error() string {
	return fmt.Sprintf("(error) ERR Protocol error: %s", p.Message)
}

type ScannerError struct {
	Message string
}

func (s ScannerError) Error() string {
	return fmt.Sprintf("(error) ERR Input Scan error: %s", s.Message)
}
