package main

import (
	"fmt"
	"net"
)

func main() {

	ip := net.IPv4(127, 0, 0, 1)
	fmt.Printf("%s\n", ip.String())

}
