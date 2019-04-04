package connect

import (
	"git.iotserv.com/iotserv/utils/io"
	"net"
	"time"
	//"github.com/xtaci/smux"
	"crypto/tls"
	"fmt"
)

func JoinTCP(stream net.Conn, ip string, port int) error {
	c, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), time.Second*30)
	if err != nil {
		return err
	}
	go io.Join(stream, c)
	return nil
}

func JoinSTCP(stream net.Conn, ip string, port int) error {
	c, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", ip, port), nil)
	if err != nil {
		return err
	}
	go io.Join(stream, c)
	return nil
}
