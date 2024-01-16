package service

import (
	"github.com/OpenIoTHub/utils/io"
	"github.com/OpenIoTHub/utils/models"
	"log"
	"net"
	"strconv"
)

func ListenMulticastUDP(stream net.Conn, service *models.NewService) error {
	host, port, err := net.SplitHostPort(service.Config)
	if err != nil {
		return err
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	l, err := net.ListenMulticastUDP("udp4", nil, &net.UDPAddr{
		IP:   net.ParseIP(host),
		Port: portInt,
	})
	if err != nil {
		log.Println(err.Error())
		return err
	}

	go io.Join(stream, l)

	//var message = make(chan []byte, 100)
	//go func() {
	//	buf := make([]byte, 2048)
	//	for {
	//		size, _, err := l.ReadFromUDP(buf)
	//		if err != nil {
	//			log.Println(err)
	//			return
	//		}
	//		log.Println(string(buf))
	//		go func() {
	//			if size > 0 {
	//				msg := make([]byte, size)
	//				copy(msg, buf[0:size])
	//				message <- msg
	//			}
	//		}()
	//	}
	//}()
	//go func() {
	//	for {
	//		msgin := <-message
	//		_, err = stream.Write(msgin)
	//		if err != nil {
	//			return
	//		}
	//		time.Sleep(time.Second)
	//	}
	//}()
	return nil
}
