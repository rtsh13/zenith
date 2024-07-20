package client

import (
	"fmt"
	"os"
	"strings"

	"github.com/zenith"
	redis "github.com/zenith/redis-protocol"
)

type Client struct {
	p redis.Protocol
}

func NewClient(p redis.Protocol) *Client {
	return &Client{p: p}
}

func (c *Client) Exec(input []string) {
	if err := c.Validate(input); err != nil {
		os.Stderr.Write([]byte(err.Error() + "\n"))
		return
	}

	fmt.Fprint(os.Stdout, c.p.Serialize(input))
}

func (c *Client) Validate(args []string) error {
	cmd := args[0]

	count, ok := zenith.Arguments(strings.ToUpper(cmd))
	if !ok {
		return zenith.UnknownCommand{Command: cmd, Args: args[1:]}
	}

	if len(args)-1 != count {
		return zenith.InvalidArgs{Command: cmd}
	}

	return nil
}
