package main

import (
	"log"
	"sync"

	. "github.com/v-gu/mmc/common"
	"github.com/v-gu/mmc/mysql"
)

var (
	wg        sync.WaitGroup
	serverMap = &ServerMap{Map: make(map[string]*Server, 50)}
	services  = make(chan int, 100)
)

func main() {
	log.Print("Starting mmc Server...")

	// servers = append(servers, Server{Host: "127.0.0.1", Port: 3306})
	// servers = append(servers, Server{Host: "192.168.11.116", Port: 3306})
	// servers = append(servers, Server{Host: "192.168.11.116", Port: 3307})
	// server := Server{Host: "127.0.0.1", Port: struct{{1, "ok"}}
	server := NewServer("localhost", "127.0.0.1", 3306)
	serverMap.Map[server.Name] = server
	server = NewServer("local2", "192.168.11.116", 3306)
	serverMap.Map[server.Name] = server
	server = NewServer("local3", "192.168.11.116", 3307)
	serverMap.Map[server.Name] = server

	// load services
	mysql.Load(&wg, serverMap)

	wg.Wait()
}
