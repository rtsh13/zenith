package server

import (
	"fmt"
	"net"
	"os"
	"strings"

	pkg "github.com/zenith"
	resp "github.com/zenith/redis-protocol"
)

type Server interface {
	Listen()
}

type server struct {
	protocol resp.Protocol
	db       dbOps
	listener net.Listener
}

func New() Server {
	protocol := resp.New()
	d := newDatabase()

	// PORT should be fetched from cfg
	// if not found, assign default port
	// timeout should follow the port pattern
	conn, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error : %v in initialising TCP connection", err.Error())
		os.Exit(1)
	}

	return &server{protocol: protocol, db: d, listener: conn}
}

type respond func(string, error) string

func dbSerializer(protocol resp.Protocol) respond {
	return func(response string, err error) string {
		if err != nil {
			return protocol.Serialize([]string{err.Error()})
		}

		return protocol.Serialize([]string{response})
	}
}

func (s *server) exec(input string) string {
	responder := dbSerializer(s.protocol)

	instructions, err := s.protocol.Deserialize(input)
	if err != nil {
		return responder("", err)
	}

	cmd := strings.Split(instructions.String(), " ")

	switch strings.ToUpper(cmd[0]) {
	case pkg.SetCMD:
		s.db.Set(cmd[1], cmd[2])
		return responder(pkg.OK, nil)
	case pkg.GetCMD:
		return responder(s.db.Get(cmd[1]), nil)
	case pkg.DelCMD:
		s.db.Delete(cmd[1])
		return responder(pkg.OK, nil)
	case pkg.EchoCMD:
		return responder(s.db.Echo(cmd[1]), nil)
	case pkg.PingCMD:
		return responder(s.db.Ping(), nil)
	default:
		return responder("", pkg.UnknownCommand{Command: cmd[0], Args: cmd[1:]})
	}
}
