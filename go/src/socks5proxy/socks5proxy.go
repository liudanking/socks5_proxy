package socks5proxy

import (
	"fmt"
	"io"
	"log"
	"net"
	"secureconn"
	"socks5"
	"strconv"
)

var (
	DefaultKey = []byte{
		102, 57, 31, 13, 11, 131, 64, 191,
		211, 221, 171, 121, 176, 173, 205, 1,
		61, 5, 3, 7, 19, 23, 41, 37,
		53, 61, 71, 91, 83, 99, 100}
)

const (
	DefaultEncType     = secureconn.RC4
	DefaultClienProxy  = ":1081"
	DefaultServerProxy = "127.0.0.1:55467"
)

type Socks5Proxy struct {
	key     []byte
	encType int
}

func (s *Socks5Proxy) handleConnect(conn net.Conn, isClient bool, proxy string) {
	buf := make([]byte, 262, 262)
	if _, err := io.ReadFull(conn, buf[:2]); err != nil {
		log.Fatal(err)
	}

	// 1. version
	if buf[0] != 0x05 {
		fmt.Printf("version 0x%02x not support", buf[0])
	}

	length := int(buf[1])
	io.ReadFull(conn, buf[:length])
	fmt.Println(buf[:length])
	conn.Write([]byte{0x05, 0x00})

	// 2. request
	io.ReadFull(conn, buf[:4])
	fmt.Println(buf[:4])
	cmd := int(buf[1])
	addrType := int(buf[3])
	addr := ""
	if addrType == 1 { // ipv4
		io.ReadFull(conn, buf[:4])
		addr = net.IPv4(buf[0],
			buf[1],
			buf[2],
			buf[3]).String()
		fmt.Println("address:", addr, buf[:4])
	} else if addrType == 3 { // domain name
		io.ReadFull(conn, buf[:1])
		length := int(buf[0])
		io.ReadFull(conn, buf[:length])
		addr = string(buf[:length])
		fmt.Println("domain: ", addr)
	} else if addrType == 4 { // ip v6
		fmt.Println("ipv6 address not support")
	} else {
		fmt.Println("address type not support: ", addrType)
	}
	io.ReadFull(conn, buf[:2])
	port := int(buf[0])<<8 + int(buf[1])
	reply := []byte{0x05, 0x00, 0x00, 0x01}

	var remote net.Conn
	if cmd == 0x01 { //   0x01: tcp connection
		addrDest := fmt.Sprintf("%s:%d", addr, port)
		var remoteTmp net.Conn
		var err error
		if isClient {
			remoteTmp, err = socks5.DialSocks5(proxy, addrDest, s.encType, s.key)
		} else {
			remoteTmp, err = net.Dial("tcp", addrDest)
		}
		if err != nil {
			fmt.Println("locate: ", err)
			conn.Close()
			//remoteTmp.Close()
			return
		}
		fmt.Println("connected to ", addr, port)
		remote = remoteTmp
		host, portStr, _ := net.SplitHostPort(remote.RemoteAddr().String())
		remoteIP := net.ParseIP(host)
		reply = append(reply, remoteIP[0], remoteIP[1], remoteIP[2], remoteIP[3])
		port, _ := strconv.ParseUint(portStr, 10, 16)
		reply = append(reply, byte(port>>8), byte(port))
	} else if cmd == 0x02 { //	0x02: tcp bind
		fmt.Println("0x02 BIND")
		return
	} else if cmd == 0x03 { //	0x03: udp associate
		fmt.Println("0X03 UDP ACCOCIATE")
		return
	} else { // command not supported
		reply = []byte{0x05, 0x07, 0x00, 0x01}
		fmt.Printf("cmd:%d, %s:%d", cmd, addr, port)
	}
	conn.Write(reply)

	if reply[1] == 0x00 { // transfer data
		if cmd == 1 {
			go handleTCP(conn, remote)
			go handleTCP(remote, conn)
		}
	}

}

func handleTCP(from, to net.Conn) {
	buf := make([]byte, 4096, 4096)
	for {
		length := 0
		if _length, _ := from.Read(buf); /*err == io.EOF &&*/ _length == 0 {
			length = _length
			//fmt.Println("io read over", length, err)
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
