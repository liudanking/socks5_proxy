package socks5proxy

import (
	"fmt"
	"log"
	"net"
	"secureconn"
)

type Socks5ProxyServer struct {
	Socks5Proxy
}

func (s *Socks5ProxyServer) ListenAndServe(network, localAddr string, encType int, key []byte) {
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
		}
		sConn := secureconn.MakeSecureConn(conn, s.encType, s.key)
		go s.handleConnect(sConn, false, "")
	}
}
