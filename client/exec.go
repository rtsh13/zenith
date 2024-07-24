package client

import (
	"fmt"
	"os"
	"strings"

	"github.com/zenith"
	redis "github.com/zenith/redis-protocol"
	"github.com/zenith/server"
)

type Client struct {
	p redis.Protocol
	s server.Server
}

func NewClient(p redis.Protocol) *Client {
	s := server.New()
	return &Client{p: p, s: s}
}

func (c *Client) Exec(input []string) {
	if err := c.Validate(input); err != nil {
		os.Stderr.Write([]byte(err.Error() + "\n"))
		return
	}

	response := c.s.Exec(c.p.Serialize(input))

	clientOutput, err := c.p.Deserialize(response)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}

	fmt.Fprint(os.Stdout, clientOutput.String())
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
