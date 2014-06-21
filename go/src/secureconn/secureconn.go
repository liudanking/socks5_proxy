// Copyright 2014, liudanking. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package secureconn

import (
	"crypto/rc4"
	"fmt"
	"net"
)

const (
	RC4  = iota // currently, only support RC4
	PASS        // do not encrypt te wire
)

type SecureConn struct {
	net.Conn
	key         []byte
	encType     int
	cipherRead  interface{}
	cipherWrite interface{}
}

func MakeSecureConn(conn net.Conn, encType int, key []byte) (sConn SecureConn) {
	sConn = SecureConn{}
	sConn.Conn = conn
	sConn.encType = encType
	sConn.key = key
	(&sConn).buildCipher(encType, key)
	return
}

func NewSecureConn(conn net.Conn, encType int, key []byte) (sConn *SecureConn) {
	sConn = &SecureConn{}
	sConn.Conn = conn
	sConn.encType = encType
	sConn.key = key
	sConn.buildCipher(encType, key)
	return
}

func DialSecureConn(network, address string, encType int, key []byte) (sConn SecureConn, err error) {
	sConn.encType = encType
	sConn.key = key
	(&sConn).buildCipher(encType, key)
	sConn.Conn, err = net.Dial(network, address)
	return
}

func (c *SecureConn) buildCipher(encType int, key []byte) {
	switch encType {
	case RC4:
		c.cipherRead, _ = rc4.NewCipher(key)
		c.cipherWrite, _ = rc4.NewCipher(key)
	case PASS:
		// do notheing, just pass
	default:
		fmt.Println("enctype invalid")
	}
}

// override read function
func (c SecureConn) Read(b []byte) (length int, err error) {
	length, err = c.Conn.Read(b)
	if length > 0 && c.encType != PASS {
		dst := make([]byte, length, length)
		c.decrypt(dst, b[:length])
		copy(b, dst)
	}
	return length, err
}

// override wirte function
func (c SecureConn) Write(b []byte) (length int, err error) {
	if c.encType != PASS {
		dst := make([]byte, len(b), len(b))
		copy(dst, b)
		c.encrypt(dst, b)
		length, err = c.Conn.Write(dst)
	} else {
		length, err = c.Conn.Write(b)
	}
	return length, err
}

func (c SecureConn) encrypt(dst, src []byte) {
	if c.encType == RC4 {
		if cipher, ok := c.cipherWrite.(*rc4.Cipher); ok {
			cipher.XORKeyStream(dst, src)
			//fmt.Println("ENC")
		}
	} else if c.encType == PASS {
		// just pass do nothing
	} else {
		fmt.Println("encType not support: ", c.encType)
	}
}

func (c SecureConn) decrypt(dst, src []byte) {
	if c.encType == RC4 {
		if cipher, ok := c.cipherRead.(*rc4.Cipher); ok {
			cipher.XORKeyStream(dst, src)
			//fmt.Println("DEC")
		}
	} else if c.encType == PASS {
		// just pass do nothing
	} else {
		fmt.Println("encType not support: ", c.encType)
	}
}

func Prints() {
	fmt.Println("...")
}
