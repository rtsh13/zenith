package client

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/zenith"

	resp "github.com/zenith/redis-protocol"
)

type client struct {
	protocol resp.Protocol
}

func New() *client {
	return &client{protocol: resp.New()}
}

func (c *client) Exec(input []string) {
	if err := validate(input); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	conn, err := net.Dial("tcp", ":6379")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	_, err = conn.Write([]byte(c.protocol.Serialize(input)))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
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
