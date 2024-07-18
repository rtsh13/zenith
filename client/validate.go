package client

import (
	"strings"

	"github.com/zenith"
)

func (c *CLI) Validate() error {
	cmd := c.Args[0]

	argsCount, ok := zenith.AllowedArgs(strings.ToUpper(cmd))
	if !ok {
		return zenith.UnknownCommand{Command: cmd, Args: c.Args[1:]}
	}

	if len(c.Args)-1 != argsCount {
		return zenith.InvalidArgs{Command: cmd}
	}

	return nil
}
