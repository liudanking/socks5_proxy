package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:1087")
	if err != nil {
		fmt.Println(err)
		return
	}

	message := "hello udp"

	b := []byte(message)
	for {
		conn.Write(b)
		time.Sleep(time.Second)
	}
}
