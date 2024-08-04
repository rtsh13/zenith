package server

import (
	"fmt"
	"os"
	"strings"

	pkg "github.com/zenith"
)

func (s *server) set(args ...string) string {
	key := args[1]
	value := args[2]

	fmt.Fprintf(os.Stdout, "query : [%v], args : [%v]", args[0], strings.Join(args[1:], ","))
	s.db.SET(key, value)

	return s.dbSerializer()(pkg.OK)
}

func (s *server) get(args ...string) string {
	fmt.Fprintf(os.Stdout, "query : [%v], args : [%v]", args[0], strings.Join(args[1:], ","))

	if len(args) > 2 {
		return s.dbSerializer()(s.db.MGET(args[1:]...))
	}

	return s.dbSerializer()(s.db.GET(args[1]))
}

func (s *server) delete(args ...string) string {
	fmt.Fprintf(os.Stdout, "query : [%v], args : [%v]", args[0], strings.Join(args[1:], ","))

	s.db.DELETE(args[1])
	return s.dbSerializer()((pkg.OK))
}

func (s *server) ping(args ...string) string {
	fmt.Fprintf(os.Stdout, "query : [%v]", args[0])
	return s.dbSerializer()(s.db.PING())
}

func (s *server) echo(args ...string) string {
	fmt.Fprintf(os.Stdout, "query : [%v], args : [%v]", args[0], strings.Join(args[1:], ","))
	return s.dbSerializer()(s.db.ECHO(args[1:]...))
}

func (s *server) incr(args ...string) string {
	fmt.Fprintf(os.Stdout, "query : [%v], args : [%v]", args[0], strings.Join(args[1:], ","))
	return s.dbSerializer()(s.db.INCR(args[1]))
}
