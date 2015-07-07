package main

import (
	"flag"
	"fmt"
	"github.com/liudanking/socks5_proxy/secureconn"
	"github.com/liudanking/socks5_proxy/socks5proxy"
)

func main() {
	mode := flag.String("m", "server", "Running Mode: server or client")
	port := flag.Int("p", 1081, "Listening port")
	key := flag.String("k", "", "Encryption key")
	host := flag.String("h", "localhost:1080", "Remote host address")
	flag.Parse()

	fmt.Println(*mode, *port, *key, *host)

	encType := secureconn.RC4
	if *key == "" {
		encType = secureconn.PASS
	}

	if *mode == "server" {
		sps := &socks5proxy.Socks5ProxyServer{}
		sps.ListenAndServe("tcp", fmt.Sprintf(":%d", *port), encType, []byte(*key))
	} else {
		spc := &socks5proxy.Socks5ProxyClient{}
		spc.ListenAndServe("tcp", fmt.Sprintf(":%d", *port), *host, encType, []byte(*key))
	}

}
