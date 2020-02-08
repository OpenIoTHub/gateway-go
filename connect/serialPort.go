package connect

import (
	"fmt"
	"github.com/OpenIoTHub/utils/io"
	"github.com/jacobsa/go-serial/serial"
	"net"
)

func JoinSerialPort(stream net.Conn, options serial.OpenOptions) error {
	conn, err := serial.Open(options)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	go io.Join(stream, conn)
	return nil
}
