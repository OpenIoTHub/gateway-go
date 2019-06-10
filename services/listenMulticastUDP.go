package services

import (
	"git.iotserv.com/iotserv/utils/models"
	"iotserv/utils/io"
	"net"
	"strconv"
)

func listenMulticastUDP(stream net.Conn, service *models.NewService) error {
	host, port, err := net.SplitHostPort(service.Config)
	if err != nil {
		return err
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return err
	}
	l, err := net.ListenMulticastUDP("udp4", nil, &net.UDPAddr{
		IP:   net.ParseIP(host),
		Port: portInt,
	})
	if err != nil {
		return err
	}
	go io.Join(stream, l)
	return nil
}
