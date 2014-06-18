package main

import (
	"fmt"
	"net"
)

func main() {
	ip := net.ParseIP("127.0.0.1")
	//ipBytes, _ := ip.MarshalText()
	for _, value := range ip[len(ip)] {
		fmt.Printf("0x%02x ", value)
	}
}
