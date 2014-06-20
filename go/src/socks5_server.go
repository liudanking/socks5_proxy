package main

import (
	_ "bufio"
	_ "errors"
	"fmt"
	"io"
	"log"
	"net"
	"secureconn"
	"strconv"
)

type Socks5ProxyServer struct {
	key     []byte
	encType int
}

func (sp *Socks5ProxyServer) ListenAndServe(network, localAddr string) {
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

func handle(conn *secureconn.SecureConn) {
	buf := make([]byte, 262, 262)
	if _, err := io.ReadFull(conn, buf[:3]); err != nil {
		log.Fatal(err)
	}

	// 1. version
	if buf[0] != 0x05 {
		fmt.Printf("version 0x%02x not support", buf[0])
	}
	conn.Write([]byte{0x05, 0x00})
	// 2. request
	io.ReadFull(conn, buf[:4])
	cmd := int(buf[1])
	addrType := int(buf[3])
	addr := ""
	if addrType == 1 { // ipv4
		io.ReadFull(conn, buf[:4])
		addr = net.IPv4(buf[0], buf[1], buf[2], buf[3]).String()
		fmt.Println("address:", addr, buf[:4])
	} else if addrType == 3 { // domain name
		io.ReadFull(conn, buf[:1])
		length := int(buf[0])
		io.ReadFull(conn, buf[:length])
		addr = string(buf[:length])
	}
	io.ReadFull(conn, buf[:2])
	port := int(buf[0])<<8 + int(buf[1])
	reply := []byte{0x05, 0x00, 0x00, 0x01}

	var remote net.Conn
	if cmd == 1 { //   1. tcp connection
		addrDest := fmt.Sprintf("%s:%d", addr, port)
		remoteTmp, err := net.Dial("tcp", addrDest)
		if err != nil {
			fmt.Println("locate: ", err)
			return
		}

		fmt.Println("connected to ", addr, port)
		remote = remoteTmp
		host, portStr, _ := net.SplitHostPort(remote.RemoteAddr().String())
		remoteIP := net.ParseIP(host)
		reply = append(reply, remoteIP[0], remoteIP[1], remoteIP[2], remoteIP[3])
		port, _ := strconv.ParseUint(portStr, 10, 16)
		reply = append(reply, byte(port>>8), byte(port))
	} else { // command not supported
		reply = []byte{0x05, 0x07, 0x00, 0x01}
	}

	conn.Write(reply)

	if reply[1] == 0x00 { // transfer data
		if cmd == 1 {
			if false { // only for development
				go handle_connection_direct(conn, remote)
				go handle_connection_direct(remote, conn)
			} else {
				go handle_connection_encrypt(conn, remote)
				go handle_connection_encrypt(remote, conn)
			}
		}
	}

}

func handle_connection_direct(from, to net.Conn) {
	io.Copy(to, from)
	from.Close()
	to.Close()
}

func handle_connection_encrypt(from, to net.Conn) {
	buf := make([]byte, 4096, 4096)
	for {
		length := 0
		if _length, err := from.Read(buf); /*err == io.EOF &&*/ _length == 0 {
			length = _length
			fmt.Println("io read over", length, err)
			break
		} else {
			length = _length
		}
		to.Write(buf[:length])
		//fmt.Printf("length: %d\n", length)
	}
	from.Close()
	to.Close()
}

func main() {
	sp := &Socks5ProxyServer{}
	sp.ListenAndServe("tcp", ":1080")

}
