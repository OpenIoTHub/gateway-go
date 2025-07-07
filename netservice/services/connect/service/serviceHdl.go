package service

import (
	"github.com/OpenIoTHub/gateway-go/v1/netservice/services/connect/service/mdns"
	"github.com/OpenIoTHub/gateway-go/v1/netservice/services/connect/tapTun"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	"net"
)

func ServiceHdl(stream net.Conn, service *models.NewService) error {
	switch service.Type {
	case "tap":
		return tapTun.NewTap(stream, service)
	case "tun":
		return tapTun.NewTun(stream, service)
	case "mDNSFind":
		return mdns.MdnsManager.FindAllmDNS(stream, service)
	case "scanPort":
		return ScanPort(stream, service)
	case "ListenMulticastUDP":
		return ListenMulticastUDP(stream, service)
	case "GetSystemStatus":
		return GetSystemStatus(stream, service)
	default:
		err := msg.WriteMsg(stream, &models.JsonResponse{Code: 1, Msg: "Failed", Result: "Unknown service type"})
		if err != nil {
			return err
		}
		return stream.Close()
	}
}
