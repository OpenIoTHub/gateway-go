package connect

import (
	"github.com/OpenIoTHub/utils/io"
	"net"
	"strconv"
	"time"
	//"github.com/xtaci/smux"
	"crypto/tls"
)

func JoinTCP(stream net.Conn, ip string, port int) error {
	c, err := net.DialTimeout("tcp", net.JoinHostPort(ip, strconv.Itoa(port)), time.Second*30)
	if err != nil {
		return err
	}
	go io.Join(stream, c)
	return nil
}

func JoinSTCP(stream net.Conn, ip string, port int) error {
	c, err := tls.Dial("tcp", net.JoinHostPort(ip, strconv.Itoa(port)), nil)
	if err != nil {
		return err
	}
	go io.Join(stream, c)
	return nil
}
