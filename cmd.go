package zenith

const (
	SET  = "SET"
	GET  = "GET"
	DEL  = "DEL"
	ECHO = "ECHO"
	PING = "PING"
	MGET = "MGET"
)

var commandArgs = map[string]int{
	SET:  2,
	GET:  1,
	DEL:  1,
	ECHO: 1,
	PING: 0,
	MGET: 100,
}

func Arguments(cmd string) (int, bool) {
	count, isValid := commandArgs[cmd]
	return count, isValid
}

const (
	Carraige = "\r"
	LineFeed = "\n"
	CRLF     = Carraige + LineFeed

	OK  = "OK"
	ERR = "-ERR"
)
