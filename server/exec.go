package server

import (
	"strings"

	pkg "github.com/zenith"
	resp "github.com/zenith/redis-protocol"
)

type Server interface {
	Exec(input string) string
}

type server struct {
	protocol resp.Protocol
	Database dbOps
}

func New() Server {
	protocol := resp.New()
	d := newDatabase()

	return &server{protocol: protocol, Database: d}
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

func (s *server) Exec(input string) string {
	responder := dbSerializer(s.protocol)

	instructions, err := s.protocol.Deserialize(input)
	if err != nil {
		return responder("", err)
	}

	cmd := strings.Split(instructions.String(), " ")

	switch cmd[0] {
	case pkg.SetCMD:
		s.Database.Set(cmd[1], cmd[2])
		return responder(pkg.OK, nil)
	case pkg.GetCMD:
		return responder(s.Database.Get(cmd[1]), nil)
	case pkg.DelCMD:
		s.Database.Delete(cmd[1])
		return responder(pkg.OK, nil)
	case pkg.EchoCMD:
		return responder(s.Database.Echo(cmd[1]), nil)
	case pkg.PingCMD:
		return responder(s.Database.Ping(), nil)
	default:
		return responder("", pkg.UnknownCommand{Command: cmd[0], Args: cmd[1:]})
	}
}
