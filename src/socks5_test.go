package main

import (
	"secureconn"
	"socks5"
	"testing"
)

func Test_DialSocks5(t *testing.T) {
	if _, err := socks5.DialSocks5("127.0.0.1:1080", "www.baidu.com:80", secureconn.RC4, []byte{1, 2, 3}); err != nil {
		t.Error(err)
	}
}
