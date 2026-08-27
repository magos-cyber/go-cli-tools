package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {
	host := flag.String("host", "127.0.0.1", "Host to scan")
	startPort := flag.Int("start", 1, "Start port")
	endPort := flag.Int("end", 1024, "End port")
	flag.Parse()

	var wg sync.WaitGroup
	for port := *startPort; port <= *endPort; port++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			addr := fmt.Sprintf("%s:%d", *host, p)
			conn, err := net.DialTimeout("tcp", addr, 500*time.Millisecond)
			if err == nil {
				conn.Close()
				fmt.Printf("Port %d: OPEN\n", p)
			}
		}(port)
	}
	wg.Wait()
}
