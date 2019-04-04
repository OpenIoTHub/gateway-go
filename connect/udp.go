package connect

import (
	"git.iotserv.com/iotserv/utils/io"
	"net"
	"time"
	//"github.com/xtaci/smux"
	"fmt"
)

func JoinUDP(stream net.Conn, ip string, port int) error {
	c, err := net.DialTimeout("udp", fmt.Sprintf("%s:%d", ip, port), time.Second*30)
	if err != nil {
		return err
	}
	go io.Join(stream, c)
	return nil
}
