package zenith

const (
	SetCMD  = "SET"
	GetCMD  = "GET"
	DelCMD  = "DEL"
	EchoCMD = "ECHO"
	PingCMD = "PING"
)

var cmdArgCounts = map[string]int{
	SetCMD:  2,
	GetCMD:  1,
	DelCMD:  1,
	EchoCMD: 1,
	PingCMD: 0,
}

func AllowedArgs(cmd string) (int, bool) {
	count, isValid := cmdArgCounts[cmd]
	return count, isValid
}
