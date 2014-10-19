package main

import (
	"fmt"
	"net"
)

var updClient = make(map[string]*net.UDPAddr)

func main() {

	addr, _ := net.ResolveUDPAddr("udp", ":0")
	conn, err := net.ListenUDP("udp", addr)
	fmt.Println(conn.LocalAddr())
	if err != nil {
		fmt.Println(err)
		return
	}

	b := make([]byte, 4096, 4096)
	for {
		n, addr, _ := conn.ReadFromUDP(b)
		fmt.Println(addr)
		fmt.Println(string(b[:n]))
	}
}
