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
	DefaultBindServer = []byte{106, 186, 114, 228}
)

const (
	DefaultEncType     = secureconn.RC4
	DefaultClienProxy  = ":1081"
	DefaultServerProxy = "127.0.0.1:55467"
)

type Socks5Proxy struct {
	key        []byte
	encType    int
	bindServer string // "2.3.4.5"
	udpServer  string // "2.3.4.5:2014"
}

func (s *Socks5Proxy) startUdpServer(net, addr string) {
	udpAddr, _ := net.ResolveUDPAddr(net, addr)
	conn, err := net.ListenUDP(net, udpAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	b := make([]byte, 4096, 4096)
	n, remoteAddr, errRead := conn.ReadFromUDP(b)
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
	reply := []byte{0x05, 0x00, 0x00, 0x01} // todo: handle ipv4 and ipv6

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
		remoteIP := net.ParseIP(host).To4()
		reply = append(reply, remoteIP[0], remoteIP[1], remoteIP[2], remoteIP[3])
		fmt.Println("remote ip: ", remoteIP[:4])
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

func (s *Socks5Proxy) serverHandleConnect(conn net.Conn) {
	buf := make([]byte, 262, 262)
	if _, err := io.ReadFull(conn, buf[:2]); err != nil {
		log.Fatal(err)
	}

	// 1. version
	if buf[0] != 0x05 {
		fmt.Printf("version 0x%02x not support", buf[0])
		conn.Close()
		return
	}

	length := int(buf[1])
	io.ReadFull(conn, buf[:length])
	conn.Write([]byte{0x05, 0x00}) // only support no auth method

	// 2. request
	cmd, addrType, addrStr, _, port := s.parseRequest(conn)

	switch cmd {
	case 0x01: // tcp connect
		fmt.Println("CONNECT")
		s.cmdConnect(conn, addrType, addrStr, port)
	case 0x02: // tcp bind
		fmt.Println("BIND")
		expectDst := fmt.Sprintf("%s:%d", addrStr, port)
		s.cmdBind(conn, expectDst)
	case 0x03: // udp associate
		fmt.Println("UDP ASSOCIATE")
	default: // unsupport cmd
		fmt.Printf("cmd 0x%02x not supported\n", cmd)
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

func (s *Socks5Proxy) parseRequest(conn net.Conn) (cmd byte, addrType byte, addrStr string, addrByte []byte, port int) {
	buf := make([]byte, 262, 262)

	io.ReadFull(conn, buf[:4])
	cmd = buf[1]
	addrType = buf[3]
	length := 0
	if addrType == 0x01 { // ipv4
		length = 4
		io.ReadFull(conn, buf[:4])
		addrStr = net.IPv4(buf[0],
			buf[1],
			buf[2],
			buf[3]).String()
		fmt.Println("ipv4 address:", addrStr)
	} else if addrType == 0x03 { // domain name
		io.ReadFull(conn, buf[:1])
		length = int(buf[0])
		io.ReadFull(conn, buf[:length])
		addrStr = string(buf[:length])
		fmt.Println("domain: ", addrStr)
	} else if addrType == 0x04 { // ip v6
		length = 16
		io.ReadFull(conn, buf[:16])
		net.IP(buf[:16]).String()
		fmt.Println("ipv6 address:", buf[:16])
	} else {
		fmt.Println("address type not support: ", addrType)
	}
	addrByte = make([]byte, length, length)
	copy(addrByte, buf[:length])

	io.ReadFull(conn, buf[:2])
	port = int(buf[0])<<8 + int(buf[1])

	return
}

func (s *Socks5Proxy) cmdConnect(conn net.Conn, addrType byte, dstAddr string, dstPort int) {
	reply := []byte{0x05, 0x00, 0x00} // todo, only support ipv4 now
	addrDest := ""
	if addrType == 0x04 { // ipv6
		addrDest = fmt.Sprintf("[%s]:%d", dstAddr, dstPort)
	} else {
		addrDest = fmt.Sprintf("%s:%d", dstAddr, dstPort)
	}
	fmt.Println("prepare connect to ", addrDest)
	remote, err := net.Dial("tcp", addrDest)
	if err != nil {
		fmt.Println("locate: ", err)
		conn.Close()
		return
	}
	fmt.Println("connected to ", addrDest, remote.RemoteAddr().String())
	host, portStr, _ := net.SplitHostPort(remote.RemoteAddr().String())

	remoteIP := net.ParseIP(host).To4()
	if remoteIP == nil {
		remoteIP = net.ParseIP(host).To16()
		if remoteIP != nil {
			reply = append(reply, 0x04) // ipv6
		} else {
			log.Fatal("remote address not ipv4/ipv6, is :", host)
		}
	} else {
		reply = append(reply, 0x01) // ipv4
	}
	for i := 0; i < len(remoteIP); i++ {
		reply = append(reply, remoteIP[i])
	}
	port, _ := strconv.ParseUint(portStr, 10, 16)
	reply = append(reply, byte(port>>8), byte(port))
	fmt.Println("reply bytes:", reply)

	// write to client
	conn.Write(reply)
	go handleTCP(conn, remote)
	go handleTCP(remote, conn)
}

func (s *Socks5Proxy) cmdBind(conn net.Conn, expectDst string) {
	listener, err := net.Listen("tcp", ":0") // auto assign a port
	if err != nil {
		fmt.Println("CMD BIND ERROR, listen failed. ", err)
		conn.Close()
		return
	}

	retBytes := make([]byte, 262, 262)

	retBytes[0] = 0x05
	retBytes[1] = 0x00
	retBytes[2] = 0x00
	retBytes[3] = 0x01 // ipv4 bind server
	copy(retBytes[4:], net.ParseIP(s.bindServer).To4())
	_, portStr, _ := net.SplitHostPort(listener.Addr().String())
	port, _ := strconv.ParseUint(portStr, 10, 16)
	retBytes[8] = byte(port >> 8)
	retBytes[9] = byte(port)
	conn.Write(retBytes[:10])

	if remote, err := listener.Accept(); err != nil {
		fmt.Println("CMD BIND ERROR, accept failed", err)
	} else {
		fmt.Println("CMD BIND, expect dst: actual dst,", expectDst, remote.RemoteAddr().String())

		retBytes[0] = 0x05
		retBytes[1] = 0x00
		retBytes[2] = 0x00
		host, portStr, _ := net.SplitHostPort(remote.RemoteAddr().String())
		remoteIP := net.ParseIP(host).To4()
		if remoteIP == nil {
			remoteIP = net.ParseIP(host).To16()
			if remoteIP != nil {
				retBytes[3] = 0x04 // ipv6
			} else {
				log.Fatal("remote address not ipv4/ipv6, is :", host)
			}
		} else {
			retBytes[3] = 0x01 // ipv4
		}
		copy(retBytes[4:], remoteIP)
		portRemote, _ := strconv.ParseUint(portStr, 10, 16)
		retBytes[4+len(remoteIP)] = byte(port >> 8)
		retBytes[5+len(remoteIP)] = byte(portRemote)
		fmt.Println("reply bytes:", retBytes[:6+len(remoteIP)])
		remote.Write(retBytes[:6+len(remoteIP)])
		go handleTCP(conn, remote)
		go handleTCP(remote, conn)
	}

	return
}

func cmdUdpAssociate(conn net.Conn, expectDst string) {

}
