package services

import (
	"git.iotserv.com/iotserv/client/connect/tapTun"
	"git.iotserv.com/iotserv/utils/models"
	"net"
)

func serviceHdl(stream net.Conn, service *models.NewService) error {
	switch service.Type {
	case "tap":
		err := tapTun.NewTap(stream, service)
		return err
	case "tun":
		err := tapTun.NewTun(stream, service)
		return err
	case "mDNSFind":
		err := findAllmDNS(stream, service)
		return err
	case "scanPort":
		err := scanPort(stream, service)
		return err
	case "ListenMulticastUDP":
		err := listenMulticastUDP(stream, service)
		return err
	default:

	}
	return nil
}
