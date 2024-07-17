package errors

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
