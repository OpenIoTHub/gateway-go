package service

import (
	"github.com/OpenIoTHub/gateway-go/netservice/services/connect/tapTun"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	"net"
)

func ServiceHdl(stream net.Conn, service *models.NewService) error {
	switch service.Type {
	case "tap":
		err := tapTun.NewTap(stream, service)
		return err
	case "tun":
		err := tapTun.NewTun(stream, service)
		return err
	case "mDNSFind":
		err := FindAllmDNS(stream, service)
		//stream.Close()
		return err
	case "scanPort":
		err := ScanPort(stream, service)
		//stream.Close()
		return err
	case "ListenMulticastUDP":
		err := ListenMulticastUDP(stream, service)
		return err
	case "GetSystemStatus":
		err := GetSystemStatus(stream, service)
		return err
	default:
		msg.WriteMsg(stream, &models.JsonResponse{Code: 1, Msg: "Failed", Result: "Unknown service type"})
		stream.Close()
	}
	return nil
}
