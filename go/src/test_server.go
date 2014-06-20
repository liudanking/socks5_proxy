package main

import (
	_ "bufio"
	_ "errors"
	"fmt"
	"io"
	"log"
	"net"
	"secureconn"
	_ "strconv"
)

type socks5proxy struct {
}

func (sp *socks5proxy) ListenAndServe(network, localAddr string) {
	listener, err := net.Listen(network, localAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	fmt.Println("listen: ", network, localAddr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		//defer conn.Close()
		sConn := secureconn.NewSecureConn(conn, secureconn.RC4, []byte{1, 2, 3})
		go handle(sConn)
	}
}

func handle(conn secureconn.SecureConn) {
	buf := make([]byte, 4, 4)
	for {
		length, err := io.ReadFull(conn, buf[:1])
		if length != 0 {
			for _, value := range buf {
				fmt.Printf("%02x ", value)
			}
		} else {
			fmt.Println(err)
			conn.Close()
			break
		}
	}

}

func main() {
	sp := &socks5proxy{}
	sp.ListenAndServe("tcp", ":2014")

}
