package persistence

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	pkg "github.com/zenith"
)

const (
	fileName   = "appendonly.aof"
	bufferSize = 1000
)

type aof struct {
	f     *os.File
	queue chan (string)
}

type WAL interface {
	Push(string)
	File() *os.File
	Close() error
}

func New() WAL {
	aof := aof{f: nil, queue: make(chan string, bufferSize)}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	absPath := filepath.Join(wd, "../../bin", fileName)

	f, err := loader(absPath)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	aof.f = f

	t := time.NewTicker(time.Second * 1)
	go func() {
		for range t.C {
			aof.write()
		}
	}()

	return &aof
}

func createFile(path string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}

	return os.Create(path)
}

func loader(path string) (*os.File, error) {
	_, err := os.Stat(path)
	switch {
	case err == nil:
		return os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0644)
	case os.IsNotExist(err):
		return createFile(path)
	default:
		return nil, err
	}
}

func (a *aof) File() *os.File {
	return a.f
}

func (a *aof) Close() error {
	return a.f.Close()
}

func (a *aof) write() {
	for {
		select {
		case cmd, ok := <-a.queue:
			if !ok {
				fmt.Fprint(os.Stderr, "queue is prematurely closed")
				return //queue closed prematurely. ideally not to be logged
			}

			// replacing CRLF delimiter with @ to avoid appending logs to new line
			cmd = strings.Replace(cmd, pkg.CRLF, "@", -1)
			fInfo, _ := a.f.Stat()
			if fInfo.Size() != 0 {
				// if data present, move the file cursor to new line
				cmd = pkg.CRLF + cmd
			}

			_, err := a.f.WriteString(cmd)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error : [%v] in persisting query : [%v]", err, cmd)
			}
		default:
			return //no items in the queue
		}
	}
}

func (a *aof) Push(cmd string) {
	a.queue <- cmd
}
