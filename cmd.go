package zenith

const (
	SetCMD  = "SET"
	GetCMD  = "GET"
	DelCMD  = "DEL"
	EchoCMD = "ECHO"
	PingCMD = "PING"
)

var commandArgs = map[string]int{
	SetCMD:  2,
	GetCMD:  1,
	DelCMD:  1,
	EchoCMD: 1,
	PingCMD: 0,
}

func Arguments(cmd string) (int, bool) {
	count, isValid := commandArgs[cmd]
	return count, isValid
}

const (
	Carraige = "\r"
	LineFeed = "\n"
	CRLF     = Carraige + LineFeed
)
