package main

import (
	_ "bufio"
	"fmt"
	"net"
	"secureconn"
	"time"
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
	conn, err := net.Dial("tcp", "127.0.0.1:2014")
	if err != nil {
		fmt.Println(err)
	} else {
		sConn := secureconn.NewSecureConn(conn, secureconn.RC4, []byte{1, 2, 3})
		//var str string
		for {
			sConn.Write([]byte{1, 2, 3, 4})
			sConn.Write([]byte{5, 6, 7, 8})
			time.Sleep(2000 * time.Millisecond)
		}
	}
}
