package connect

import (
	"git.iotserv.com/iotserv/utils/io"
	"net"
	"strconv"
	"time"
)

func JoinUDP(stream net.Conn, ip string, port int) error {
	c, err := net.DialTimeout("udp", net.JoinHostPort(ip, strconv.Itoa(port)), time.Second*30)
	if err != nil {
		return err
	}
	go io.Join(stream, c)
	return nil
}
