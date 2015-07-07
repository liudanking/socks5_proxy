package socks5proxy

import (
	"fmt"
	"log"
	"net"
)

type Socks5ProxyClient struct {
	Socks5Proxy
}

func (s *Socks5ProxyClient) ListenAndServe(network, localAddr, proxy string, encType int, key []byte) {
	s.encType = encType
	s.key = key
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
		} else {
			fmt.Println("Accept connection: ", conn.RemoteAddr())
		}
		//defer conn.Close()
		go s.handleConnect(conn, true, proxy)
	}
}
