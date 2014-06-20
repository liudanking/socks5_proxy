package main

import (
	_ "bufio"
	_ "errors"
	"secureconn"
	"socks5proxy"
)

func main() {
	sp := &socks5proxy.Socks5ProxyServer{}

	sp.ListenAndServe("tcp", ":1080", secureconn.RC4, []byte{1, 2, 3})

}
