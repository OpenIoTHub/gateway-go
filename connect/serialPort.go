package connect

import (
	"fmt"
	"git.iotserv.com/iotserv/iotserv/utils/io"
	"github.com/jacobsa/go-serial/serial"
	"net"
)

func JoinSerialPort(stream net.Conn, options serial.OpenOptions) error {
	conn, err := serial.Open(options)
	if err != nil {
		fmt.Println("serial.Open: %v", err)
		return err
	}
	go io.Join(stream, conn)
	return nil
}
