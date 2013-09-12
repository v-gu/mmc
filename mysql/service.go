package mysql

import (
	"log"
	"net"
	"sync"
	"time"

	. "github.com/v-gu/mmc/common"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	waitGroup *sync.WaitGroup

	serverMap    *ServerMap
	portTicker   *time.Ticker = time.NewTicker(POLLING_INTERVAL)
	statusTicker *time.Ticker = time.NewTicker(POLLING_INTERVAL)
)

// Load actual server[s] to be monitored.
func Load(wg *sync.WaitGroup, sMap *ServerMap) {
	waitGroup = wg
	serverMap = sMap

	// load services
	startService(portTicker, portService)
	startService(statusTicker, statusService)
}

// Service wrapper
func startService(ticker *time.Ticker, f func(server *Server)) {
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		for {
			serverMap.RLock()
			for _, server := range serverMap.Map {
				if !CheckServerAvail(server) {
					continue
				}
				f(server)
			}
			serverMap.RUnlock()
			<-ticker.C
		}
	}()
}

// MySQL port monitoring.
func portService(server *Server) {
	server.Port.Avail = PortStat(net.ParseIP(server.Host),
		server.Port.Num)
	Debug.Printf("Server '%v:%v' status: %v\n",
		server.Host, server.Port.Num, server.Port.Avail)
}

// Mysql STATUS monitoring.
func statusService(server *Server) {
	db, err := sql.Open("mysql", "rmon:rmon@/")
	if err != nil {
		log.Print(err)
		return
	}
	rows, err := db.Query("SHOW STATUS")
	if err != nil {
		log.Print(err)
		return
	}
	for rows.Next() {
		var name, value string
		if err := rows.Scan(&name, &value); err != nil {
			log.Print(err)
			return
		}
		server.Status[name] = value
	}
	if err := rows.Err(); err != nil {
		log.Print(err)
	}
}
