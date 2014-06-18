# coding=UTF-8

import sys
import socks
import socket

print 'niuniu是笨蛋'

import socks
s = socks.socksocket()
s.setproxy(socks.PROXY_TYPE_SOCKS5, "127.0.0.1", 1080)
s.connect(("www.baidu.com",80))
s.sendall("hello")
ret = s.recv(1000)

print ret