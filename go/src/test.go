package main

import (
	"crypto/rc4"
	"fmt"
	_ "net"
)

func printHex(bytes []byte) {
	fmt.Printf("\n=======\n")
	for _, value := range bytes {
		fmt.Printf("%02x ", value)
	}
}

func main() {
	plainText := []byte{1, 2, 3, 4}
	cipher, _ := rc4.NewCipher([]byte{9, 10})
	dst := make([]byte, 4, 4)
	cipher.XORKeyStream(dst, plainText)

	printHex(plainText)
	printHex(dst)

	//
	printHex(plainText)
	dst2 := make([]byte, 4, 4)
	cipher, _ = rc4.NewCipher([]byte{9, 10})
	cipher.XORKeyStream(dst2[:2], dst[:2])
	cipher.XORKeyStream(dst2[2:4], dst[2:4])
	printHex(plainText)
	printHex(dst2)

}
