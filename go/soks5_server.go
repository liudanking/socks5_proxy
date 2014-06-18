package main

import (
	_ "bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
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
		go handle(conn)
	}
}

func handle(conn net.Conn) {
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
		addr = string(net.IPv4(buf[0], buf[1], buf[2], buf[3]))
	} else if addrType == 3 { // domain name
		io.ReadFull(conn, buf[:1])
		length := int(buf[0])
		io.ReadFull(conn, buf[:length])
		addr = string(buf[:length])
	}
	io.ReadFull(conn, buf[:2])
	port := int(buf[0])<<8 + int(buf[1])
	fmt.Printf("0x%02x 0x%02x", buf[0], buf[1])
	reply := []byte{0x05, 0x00, 0x00, 0x01}

	var remote net.Conn
	if cmd == 1 { //   1. tcp connection
		laddr := addr + ":" + fmt.Sprintf("%d", port)
		fmt.Println("tcp connect to ", laddr)
		_remote, err := net.Dial("tcp", laddr)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("already connect to ", addr, port)
		remote = _remote
		fmt.Println(remote.RemoteAddr().String())
		host, portStr, _ := net.SplitHostPort(remote.RemoteAddr().String())
		remoteIP := net.ParseIP(host)
		reply = append(reply, remoteIP[0], remoteIP[1], remoteIP[2], remoteIP[3])
		port, _ := strconv.ParseUint(portStr, 10, 16)
		reply = append(reply, byte(port>>8), byte(port))
		//append(remote, strconv.par)
	} else { // command not supported
		reply = []byte{0x05, 0x07, 0x00, 0x01}
	}
	conn.Write(reply)

	if reply[1] == 0x00 {
		if cmd == 1 {
			fmt.Println("transfer")

			go io.Copy(remote, conn)
			go io.Copy(conn, remote)
		}
	}

}

func main() {
	sp := &socks5proxy{}
	sp.ListenAndServe("tcp", ":8084")

}
