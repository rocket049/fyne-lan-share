package main

import (
	"gitee.com/rocket049/discover-go"
)

func runDiscover(ips []string) {
	server := discover.NewServer()
	for _, ip := range ips {
		server.Append("http", ip, 6868, "index", "FileServer", "Share Files")
	}

	go server.Serve(true)
}
