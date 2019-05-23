package connect

import (
	"git.iotserv.com/iotserv/utils/io"
	"net"
)

func JoinUDP(stream net.Conn, ip string, port int) error {
	addr := &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	}
	//udp还是udp4
	c, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return err
	}
	go io.Join(stream, c)
	return nil
}
