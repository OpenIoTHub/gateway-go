package service

import (
	"log"
	"net"
	"time"
	//"github.com/xtaci/smux"
	"crypto/tls"
)

// CheckTcpUdpTls 检查端口状态，如果端口可连接则状态为0，如果不可连接则状态为其他
func CheckTcpUdpTls(connType, addr string) (int, string) {
	var c net.Conn = nil
	var err error
	defer func() {
		if c != nil {
			err = c.Close()
			if err != nil {
				log.Println(err.Error())
			}
		}
	}()
	switch connType {
	case "tcp", "udp":
		c, err = net.DialTimeout(connType, addr, time.Second)
	case "tls":
		c, err = tls.Dial("tcp", addr, nil)
	default:
		return 1, "type not tcp,udp or tls"
	}
	if err != nil {
		return 1, err.Error()
	}
	return 0, ""
}
