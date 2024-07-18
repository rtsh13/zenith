package client

import (
	"fmt"
	"strings"
)

func (c *CLI) Serialize() string {
	arr := make([]string, 0)

	for _, value := range c.Args {
		arr = append(arr, fmt.Sprintf("$%d\r\n%v", len(value), value))
	}

	response := strings.Join(arr, "\r\n")

	return fmt.Sprintf("*%d\r\n%v\r\n", len(c.Args), response)
}
