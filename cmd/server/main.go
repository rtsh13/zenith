package main

import (
	"fmt"
	"os"

	"github.com/zenith/server"
)

func main() {
	s, err := server.New()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	defer s.Close()

	s.Listen()
}
