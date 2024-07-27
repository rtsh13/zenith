package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func (s *server) Listen() {
	fmt.Fprintln(os.Stdout, "listening at port :6379")

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		go s.read(conn)
	}
}

func (s *server) read(conn net.Conn) {
	reader := bufio.NewReader(conn)

	var (
		buffer = make([]byte, 1024)
		data   = strings.Builder{}
	)

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			fmt.Fprint(os.Stdout, err)
			return
		}

		data.WriteString(string(buffer[:n]))

		if n < len(buffer) {
			break
		}
	}

	response := s.exec(data.String())

	if _, err := conn.Write([]byte(response)); err != nil {
		fmt.Fprint(os.Stdout, err)
	}
}
