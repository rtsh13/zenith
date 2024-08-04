package server

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	errors "github.com/zenith/errors/server"
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
	INCR(string) any
}

func newDatabase() dbOps {
	return &database{records: make(map[string]string, 0), mu: sync.Mutex{}}
}

func (d *database) PING() string { return "PONG" }

func (d *database) SET(key, value string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.records[key] = value
}

func (d *database) GET(key string) string {
	d.mu.Lock()
	defer d.mu.Unlock()
	val, ok := d.records[key]
	if !ok {
		return "(nil)"
	}

	return val
}

func (d *database) DELETE(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.records, key)
}

func (d *database) ECHO(input ...string) []string {
	return input
}

func (d *database) MGET(keys ...string) []string {
	list := make([]string, 0)

	for _, k := range keys {
		list = append(list, d.GET(k))
	}

	return list
}

func (d *database) INCR(key string) any {
	val := d.GET(key)

	// if the key doesn't exist, seed with zero value before completing operation
	if strings.EqualFold(val, "(nil)") {
		newVal := 0
		newVal++

		d.SET(key, strconv.FormatInt(int64(newVal), 10))

		return fmt.Sprintf("(%T) %v", newVal, newVal)
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return errors.CustomError{Message: "value is not an integer or out of range"}
	}

	intVal++
	d.SET(key, strconv.FormatInt(int64(intVal), 10))

	return fmt.Sprintf("(%T) %v", intVal, intVal)
}
