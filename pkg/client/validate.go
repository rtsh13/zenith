package client

import (
	"strings"

	"github.com/zenith/pkg/command"
	"github.com/zenith/pkg/errors"
)

func (c *CLI) Validate() error {
	cmd := c.Args[0]

	argsCount, ok := command.AllowedArgs(strings.ToUpper(cmd))
	if !ok {
		return errors.UnknownCommand{Command: cmd, Args: c.Args[1:]}
	}

	if len(c.Args)-1 != argsCount {
		return errors.InvalidArgs{Command: cmd}
	}

	return nil
}
