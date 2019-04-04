package serial

import (
	"net"
)

func JoinSerialPort(stream net.Conn, port string, baud int) error {
	//c := &serial.Config{Name: port, Baud: baud}
	//s, err := serial.OpenPort(c)
	//if err != nil {
	//	fmt.Printf("serial err")
	//	return err
	//}
	//fmt.Printf("join serial")
	//go io.Join(stream, s)
	return nil
}
