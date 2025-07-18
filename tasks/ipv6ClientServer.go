package tasks

import (
	"github.com/OpenIoTHub/gateway-go/v2/chans"
	"github.com/OpenIoTHub/gateway-go/v2/config"
	"github.com/OpenIoTHub/gateway-go/v2/netservice/handle"
	"github.com/OpenIoTHub/utils/v2/models"
	"github.com/OpenIoTHub/utils/v2/msg"
	"github.com/libp2p/go-yamux"
	"log"
	"net"
	"time"
)

func RunTasks() {
	go ipv6ServerTask()
	go ipv6ClientTask()
}

// Ipv6ClientTask 接收配置创建新的Client handle
func ipv6ClientTask() {
	//	主动连接访问者的APP
	for remoteIpv6Server := range chans.ClientTaskChan {
		log.Println("====获取到任务")
		ip := remoteIpv6Server.Ipv6AddrIp
		port := remoteIpv6Server.Ipv6AddrPort
		//runId := remoteIpv6Server.RunId
		//	使用配置创建连接，并且发送带RunId的凭证给访问者
		ipv6conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
			IP:   net.ParseIP(ip),
			Port: port,
		})
		if err != nil {
			log.Println("ipv6 net.DialTCP" + err.Error())
			continue
		}
		log.Println("ipv6 net.DialTCP connected:" + remoteIpv6Server.Ipv6AddrIp)
		//TODO 发送凭证
		runIdMsg := &models.JsonResponse{}
		err = msg.WriteMsg(ipv6conn, runIdMsg)
		if err != nil {
			log.Println("ipv6 msg.WriteMsg" + err.Error())
			ipv6conn.Close()
			return
		}
		//创建session，session handle
		yamuxConfig := yamux.DefaultConfig()
		//remoteIpv6Server.EnableKeepAlive = false
		session, err := yamux.Server(ipv6conn, yamuxConfig)
		if err != nil {
			ipv6conn.Close()
			return
		}
		log.Printf("ipv6 p2p client HandleSession")
		go handle.HandleSession(session, "")
	}
}

func ipv6ServerTask() {
	listener, err := net.ListenTCP("tcp6", &net.TCPAddr{})
	if err != nil {
		log.Println(err)
		return
	}
	listenerPort := listener.Addr().(*net.TCPAddr).Port
	log.Println("ipv6 server listening on", listenerPort)
	config.Ipv6ListenTcpHandlePort = listenerPort
	//接受验证连通性，接受连接和服务请求
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		log.Println("ipv6 server handle conn", conn.RemoteAddr().String())
		//	验证token，回复
		go ipv6ClientHandle(conn)
	}
}

func ipv6ClientHandle(conn net.Conn) {
	rawMsg, err := msg.ReadMsg(conn)
	if err != nil {
		log.Println(err.Error() + "从stream读取数据错误")
		conn.Close()
		return
	}
	// TODO 验证token,RunId
	_ = rawMsg
	// Token为空
	session, err := yamux.Server(conn, yamux.DefaultConfig())
	if err != nil {
		conn.Close()
		return
	}
	log.Println("ipv6 server handle session", conn.RemoteAddr().String())
	//TODO token!!!!
	go handle.HandleSession(session, "")
}
