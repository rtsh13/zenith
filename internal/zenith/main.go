package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	z "github.com/zenith/pkg/client"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			os.Stderr.Write([]byte(fmt.Sprintf("error : [%v] reading input", err.Error())))
			continue
		}

		input = strings.TrimSpace(input)

		switch input {
		case "":
			continue
		case strings.ToLower("exit"):
			os.Exit(1)
		case strings.ToLower("help"):
			help()
		default:
			execute(input)
		}
	}
}

func help() {
	fmt.Fprintln(os.Stdout, "Usage: go run main.go <command>")
	fmt.Fprintln(os.Stdout, "Available commands:")
	fmt.Fprintln(os.Stdout, "  SET <key> <value>   : Set the value of a key")
	fmt.Fprintln(os.Stdout, "  GET <key>           : Get the value of a key")
	fmt.Fprintln(os.Stdout, "  DEL <key>           : Delete a key")
	fmt.Fprintln(os.Stdout, "  ECHO <message>      : Echo the message")
	fmt.Fprintln(os.Stdout, "  PING                : Ping the server")
	fmt.Fprintln(os.Stdout, "  help                : Show this help message")
}

func execute(input string) {
	c := z.NewClient(strings.Fields(input))
	if err := c.Validate(); err != nil {
		os.Stderr.Write([]byte(err.Error() + "\n"))
		return
	}

	fmt.Fprint(os.Stdout, c.Serialize())
}
