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

type respond func(interface{}) string

type cmdRouter func(...string) string

func (s *server) route() map[string]cmdRouter {
	return map[string]cmdRouter{
		pkg.SET:  s.set,
		pkg.GET:  s.get,
		pkg.DEL:  s.delete,
		pkg.ECHO: s.echo,
		pkg.PING: s.ping,
		pkg.MGET: s.get,
	}
}

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

func (s *server) dbSerializer() respond {
	return func(response interface{}) string {
		switch v := response.(type) {
		case string:
			return s.protocol.Serialize([]string{v})
		case []string:
			return s.protocol.Serialize(v)
		case error:
			return s.protocol.Serialize([]string{v.Error()})
		default:
			return s.protocol.Serialize([]string{pkg.ERR + ": unknown response type"})
		}
	}
}

func (s *server) exec(input string) string {
	responder := s.dbSerializer()

	instructions, err := s.protocol.Deserialize(input)
	if err != nil {
		return responder(err)
	}

	ins := strings.Split(instructions.String(), " ")

	if handler, found := s.route()[strings.ToUpper(ins[0])]; found {
		response := handler(ins...)

		if strings.EqualFold(ins[0], pkg.SET) || strings.EqualFold(ins[0], pkg.DEL) {
			s.wal.Push(input)
		}

		return response
	}

	return responder(errors.UnknownCommand{Command: ins[0], Args: ins[1:]})
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

		ins := strings.Split(instructions.String(), " ")
		if strings.EqualFold(ins[0], pkg.SET) {
			s.db.SET(ins[1], ins[2])
			continue
		}

		if strings.EqualFold(ins[0], pkg.DEL) {
			s.db.DELETE(ins[1])
			continue
		}

		mulErr.Errors = append(mulErr.Errors, fmt.Errorf("unknown command : %v to restore from AOF", ins[0]))
	}

	return mulErr
}
