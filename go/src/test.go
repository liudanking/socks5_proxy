package main

import (
	_ "bufio"
	_ "fmt"
	"secureconn"
	"socks5proxy"
)

// func f()
// {
// 	conn, err := net.Dial("tcp", "127.0.0.1:2014")
// 	if err != nil {
// 		fmt.Println(err)
// 	} else {
// 		sConn := secureconn.NewSecureConn(conn, secureconn.RC4, []byte{1, 2, 3})
// 		var str string
// 		for {
// 			fmt.Scanf("%s", str)
// 			bytes := []byte(str)
// 			for _, value := range bytes {
// 				fmt.Printf("%02x", value)
// 			}
// 			fmt.Printf("\n")
// 			sConn.Write(bytes)

// 		}
// 	}
//}

func main() {
	sp := &socks5proxy.Socks5ProxyClient{}
	sp.ListenAndServe("tcp", ":1081", "127.0.0.1:1080", secureconn.RC4, []byte{1, 2, 3})
}
