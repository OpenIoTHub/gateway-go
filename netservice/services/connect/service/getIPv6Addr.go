package service

import (
	"encoding/json"
	"github.com/OpenIoTHub/getip/v2/iputils"
	"github.com/OpenIoTHub/utils/v2/models"
	"github.com/OpenIoTHub/utils/v2/msg"
	"log"
	"net"
)

// GetIPv6Addr 传过来一个ipv6Addr+port返回一个ipv6Addr+port
func GetIPv6Addr(stream net.Conn, service *models.NewService) error {
	ipv6Map := make(map[string]interface{})

	// 获取磁盘信息
	var ipv6Addr = iputils.GetMyPublicIpv6()
	ipv6Map["Ipv6Addr"] = ipv6Addr
	rstByte, err := json.Marshal(ipv6Map)
	if err != nil {
		log.Println("json.Marshal(ipv6Map)：")
		log.Println(err.Error())
	}
	//log.Println(string(rstByte))
	err = msg.WriteMsg(stream, &models.JsonResponse{Code: 0, Msg: "Success", Result: string(rstByte)})
	if err != nil {
		stream.Close()
		log.Println("写消息错误：")
		log.Println(err.Error())
	}
	return err
}
