package service

import (
	"encoding/json"
	"github.com/OpenIoTHub/gateway-go/v2/chans"
	"github.com/OpenIoTHub/gateway-go/v2/config"
	models2 "github.com/OpenIoTHub/gateway-go/v2/models"
	"github.com/OpenIoTHub/getip/v2/iputils"
	"github.com/OpenIoTHub/utils/v2/models"
	"github.com/OpenIoTHub/utils/v2/msg"
	"log"
	"net"
)

// GetIPv6Addr 传过来一个ipv6Addr+port返回一个ipv6Addr+port，传过来的ipv6Addr+port用于这边通过channel通知连接，
// 访问端获取的ipv6Addr+port用于先校验连通性再有连接需求的时候使用连接
func GetIPv6Addr(stream net.Conn, service *models.NewService) error {
	//从service读取访问者监听ipv6端口号，这里将ipv6地址+port通过chan发送到任务创建连接并处理请求
	var remoteIpv6ServerConfig models2.Ipv6ClientHandleTask
	err := json.Unmarshal([]byte(service.Config), &remoteIpv6ServerConfig)
	if err != nil {
		log.Println("json.Unmarshal([]byte(service.Config), &config):" + err.Error())
		return err
	}
	chans.ClientTaskChan <- remoteIpv6ServerConfig
	// 获取ipv6公网地址，可能为空字符串代表没有或者没获取到
	var ipv6Addr = iputils.GetMyPublicIpv6()
	ipv6Info := models2.Ipv6ClientHandleTask{}
	//访问者只要保存提供服务的ipv6地址+端口，有连接请求时创建连接
	ipv6Info.Ipv6AddrIp = ipv6Addr
	ipv6Info.Ipv6AddrPort = config.Ipv6ListenTcpHandlePort
	rstByte, err := json.Marshal(ipv6Info)
	if err != nil {
		log.Println("json.Marshal(ipv6Map)：")
		log.Println(err.Error())
		return err
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
