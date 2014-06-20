package secureconn

import (
	"crypto/rc4"
	"fmt"
	"net"
)

const (
	RC4 = iota // currently, only support RC4
	WAKE
)

type SecureConn struct {
	net.Conn
	key         []byte
	encType     int
	cipherRead  *rc4.Cipher
	cipherWrite *rc4.Cipher
}

func NewSecureConn(conn net.Conn, encType int, key []byte) (sConn SecureConn) {
	sConn.Conn = conn
	sConn.encType = encType
	sConn.key = key
	sConn.cipherRead, _ = rc4.NewCipher(key)
	sConn.cipherWrite, _ = rc4.NewCipher(key)
	return
}

func DialSecureConn(network, address string, encType int, key []byte) (sConn SecureConn, err error) {
	sConn.encType = encType
	sConn.key = key
	sConn.cipherRead, _ = rc4.NewCipher(key)
	sConn.cipherWrite, _ = rc4.NewCipher(key)
	sConn.Conn, err = net.Dial(network, address)
	return
}

// override read function
func (c SecureConn) Read(b []byte) (int, error) {
	length, err := c.Conn.Read(b)
	if length > 0 {
		dst := make([]byte, length, length)
		c.decrypt(dst, b[:length])
		copy(b, dst)
	}
	return length, err
}

// override wirte function
func (c SecureConn) Write(b []byte) (int, error) {
	dst := make([]byte, len(b), len(b))
	copy(dst, b)
	c.encrypt(dst, b)
	length, err := c.Conn.Write(dst)
	return length, err
}

func (c SecureConn) encrypt(dst, src []byte) {
	if c.encType == RC4 {
		c.cipherWrite.XORKeyStream(dst, src)
		fmt.Println("ENC")
	} else {
		// pass
		fmt.Println("encType not support: ", c.encType)
	}
}

func (c SecureConn) decrypt(dst, src []byte) {
	if c.encType == RC4 {
		c.cipherRead.XORKeyStream(dst, src)
		fmt.Println("DEC")
	} else {
		// pass
		fmt.Println("encType not support: ", c.encType)
	}
}
