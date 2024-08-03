package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	pkg "github.com/zenith"

	errors "github.com/zenith/errors/server"
	"github.com/zenith/persistence"
	resp "github.com/zenith/redis-protocol"
)

type Server interface {
	Listen()
	Close()
}

type server struct {
	protocol resp.Protocol
	db       dbOps
	listener net.Listener
	wal      persistence.WAL
}

func New() (Server, error) {
	conn, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error : %v in initialising TCP connection", err.Error())
		return nil, err
	}

	server := server{
		listener: conn,
		protocol: resp.New(),
		db:       newDatabase(),
		wal:      persistence.New(),
	}

	if err := server.restore(); len(err.Errors) != 0 {
		return nil, err
	}

	return &server, nil
}

func (s *server) Close() {
	if err := s.listener.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "error : [%v] in terminating TCP conn", err)
	}

	if err := s.wal.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "error : [%v] in closing the  the file", err)
	}
}

type respond func(string, error) string

func (s *server) dbSerializer() respond {
	return func(response string, err error) string {
		if err != nil {
			return s.protocol.Serialize([]string{err.Error()})
		}

		return s.protocol.Serialize([]string{response})
	}
}

func (s *server) exec(input string) string {
	responder := s.dbSerializer()

	instructions, err := s.protocol.Deserialize(input)
	if err != nil {
		return responder("", err)
	}

	cmd := strings.Split(instructions.String(), " ")
	defer func(string, string) {
		if strings.EqualFold(cmd[0], pkg.SetCMD) || strings.EqualFold(cmd[0], pkg.DelCMD) {
			s.wal.Push(input)
		}
	}(input, cmd[0])

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
		return responder("", errors.UnknownCommand{Command: cmd[0], Args: cmd[1:]})
	}
}

func (s *server) restore() errors.MultipleErrors {
	mulErr := errors.MultipleErrors{}

	f := s.wal.File()
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		if err := scanner.Err(); err != nil {
			mulErr.Errors = append(mulErr.Errors, err)
			continue
		}

		if line == "" {
			continue
		}

		line = strings.ReplaceAll(line, "@", pkg.CRLF)

		instructions, err := s.protocol.Deserialize(line)
		if err != nil {
			mulErr.Errors = append(mulErr.Errors, err)
			continue
		}

		cmd := strings.Split(instructions.String(), " ")
		if cmd[0] == pkg.SetCMD {
			s.db.Set(cmd[1], cmd[2])
			continue
		}

		if cmd[0] == pkg.DelCMD {
			s.db.Delete(cmd[1])
			continue
		}

		mulErr.Errors = append(mulErr.Errors, fmt.Errorf("unknown command : %v to restore from AOF", cmd[0]))
	}

	return mulErr
}
