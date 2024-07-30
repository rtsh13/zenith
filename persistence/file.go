package persistence

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
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
	Read()
	Push(string)
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

	t := time.NewTicker(time.Second * 1)

	go func() {
		for range t.C {
			write(&aof)
		}
	}()

	aof.f = f

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
		return os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	case os.IsNotExist(err):
		return createFile(path)
	default:
		return nil, err
	}
}

func (a *aof) Read() {}

func write(a *aof) {
	for {
		select {
		case cmd, ok := <-a.queue:
			if !ok {
				fmt.Fprint(os.Stderr, "queue is prematurely closed")
				return //queue closed prematurely. ideally not to be logged
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
