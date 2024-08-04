package server

import (
	"sync"
)

type database struct {
	mu      sync.Mutex
	records map[string]string
}

type dbOps interface {
	SET(key, value string)
	GET(key string) string
	DELETE(key string)
	ECHO(input ...string) []string
	PING() string
	MGET(...string) []string
}

func newDatabase() dbOps {
	return &database{records: make(map[string]string, 0), mu: sync.Mutex{}}
}

func (d *database) Ping() string { return "PONG" }

func (d *database) Set(key, value string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.records[key] = value
}

func (d *database) Get(key string) string {
	d.mu.Lock()
	defer d.mu.Unlock()
	val, ok := d.records[key]
	if !ok {
		return "(nil)"
	}

	return val
}

func (d *database) Delete(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.records, key)
}

func (d *database) Echo(input string) string {
	return input
}

func (d *database) MGET(keys ...string) []string {
	list := make([]string, 0)

	for _, k := range keys {
		list = append(list, d.GET(k))
	}

	return list
}
