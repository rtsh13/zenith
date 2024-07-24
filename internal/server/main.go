package main

import (
	"fmt"

	"github.com/zenith/server"
)

func main() {
	s := server.New()

	fmt.Print(s.Exec("*3\r\n$3\r\nSETi\r\n$5\r\nHELLO\r\n$5\r\nWORLD\r\n "))
	fmt.Print(s.Exec("*2\r\n$3\r\nGET\r\n$5\r\nHELLO\r\n"))
	fmt.Print(s.Exec("*2\r\n$3\r\nDEL\r\n$5\r\nHELLO\r\n"))
}
