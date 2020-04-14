package connect

import (
	"github.com/OpenIoTHub/utils/io"
	"github.com/jacobsa/go-serial/serial"
	"log"
	"net"
)

func JoinSerialPort(stream net.Conn, options serial.OpenOptions) error {
	conn, err := serial.Open(options)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	go io.Join(stream, conn)
	return nil
}
