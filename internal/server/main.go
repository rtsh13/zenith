package main

import (
	"fmt"

	resp "github.com/zenith/redis-protocol"
	"github.com/zenith/server"
)

func main() {
	s := server.New()

	r := resp.New()

	response := s.Exec("+PING\r\n")
	ss, _ := r.Deserialize(response)
	fmt.Print(ss.String())

	response = s.Exec("*3\r\n$3\r\nSET\r\n$5\r\nHELLO\r\n$5\r\nWORLD\r\n ")
	ss, _ = r.Deserialize(response)
	fmt.Print(ss.String())

	response = s.Exec("*2\r\n$3\r\nGET\r\n$5\r\nHELaL\r\n")
	ss, _ = r.Deserialize(response)
	fmt.Print(ss.String())

	response = s.Exec("*2\r\n$3\r\nDEL\r\n$5\r\nHELLO\r\n")
	ss, _ = r.Deserialize(response)
	fmt.Print(ss.String())

	response = s.Exec("*2\r\n$4\r\nECHO\r\n$5\r\nHELLO\r\n")
	ss, _ = r.Deserialize(response)
	fmt.Print(ss.String())
}
