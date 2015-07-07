package secureconn

import (
	"bytes"
	"testing"
)

func Test_encrypt_decrype(t *testing.T) {
	sConn := MakeSecureConn(nil, RC4, []byte{1, 2, 3})
	src := []byte{4, 5, 6}
	dst := make([]byte, 3, 3)
	sConn.encrypt(dst, src)
	sConn = MakeSecureConn(nil, RC4, []byte{1, 2, 3})
	_src := make([]byte, 3, 3)
	sConn.decrypt(_src, dst)

	if !bytes.Equal(_src, src) {
		t.Errorf("%q != %q", _src, src)
	}
}
