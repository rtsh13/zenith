package client

import (
	"fmt"
	"net"
	"os"
	"strings"

	pkg "github.com/zenith"
	errors "github.com/zenith/errors/client"

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

	var (
		buffer = make([]byte, 1024)
		data   = strings.Builder{}
	)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		data.WriteString(string(buffer[:n]))

		if n < len(buffer) {
			break
		}
	}

	clientOutput, err := c.protocol.Deserialize(data.String())
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}

	fmt.Fprint(os.Stdout, clientOutput.String()+pkg.LineFeed)
}

func validate(args []string) error {
	cmd := args[0]

	count, ok := pkg.Arguments(strings.ToUpper(cmd))
	if !ok {
		return errors.UnknownCommand{Command: cmd, Args: args[1:]}
	}

	if strings.EqualFold(cmd, pkg.MGET) {
		if len(args) >= 1 {
			return nil
		}
	}

	if len(args)-1 != count {
		return errors.InvalidArgs{Command: cmd}
	}

	return nil
}
