package main

import (
	"strconv"
	"strings"
)

// Instance has no documentation
type Instance struct {
	Host string
	Port int
}

func instanceOf(instance string) *Instance {
	idx := strings.Index(instance, ":")
	if idx == -1 {
		return nil
	}
	host := instance[0:idx]
	var port int
	var err error
	port, err = strconv.Atoi(instance[idx+1:])
	if err != nil {
		return nil
	}
	return &Instance{Host: host, Port: port}
}
