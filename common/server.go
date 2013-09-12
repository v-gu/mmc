package common

import (
	"sync"
	"time"
)

const (
	POLLING_INTERVAL time.Duration = time.Second * 3 // 10s
)

type Server struct {
	sync.RWMutex
	Disabled bool
	Name     string
	Host     string
	Port
	Status ServerStatus
}

func NewServer(name string, host string, port int) *Server {
	server := &Server{}
	server.Name, server.Host = name, host
	server.Port.Num = port
	server.Status = make(map[string]string, 500)
	return server
}

func CheckServerAvail(server *Server) (result bool) {
	server.RLock()
	if server.Disabled {
		result = false
	} else {
		result = true
	}
	server.RUnlock()
	return result
}

type ServerMap struct {
	sync.RWMutex
	Map map[string]*Server
}

type Port struct {
	Num   int    // port number
	Avail bool   // port availability
	Err   string // error string if port is not available
}

type ServerStatus map[string]string
