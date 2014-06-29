package main

import (
	"fmt"
	"net"
)

func main() {
	addr, _ := net.ResolveUDPAddr("udp", ":1087")
	fmt.Println(addr)
	conn, err := net.ListenUDP("udp", addr)
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
