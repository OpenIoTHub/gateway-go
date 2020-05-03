package p2p

import (
	"fmt"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	nettool "github.com/OpenIoTHub/utils/net"
	"github.com/libp2p/go-yamux"
	"github.com/xtaci/kcp-go/v5"
	"log"
	"net"
	"time"
)

//var ExternalPort int
//var SendPackReqChan = make(chan *models.SendUdpPackReq,10)

func NewP2PCtrlAsServer(stream net.Conn, token *models.TokenClaims) (*yamux.Session, error) {
	if stream != nil {
		defer stream.Close()
	}
	//监听一个随机端口号，接受P2P方的连接
	externalUDPAddr, listener, err := GetP2PListener(token)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = msg.WriteMsg(stream, &models.ReqNewP2PCtrlAsClient{
		IntranetIp:   listener.LocalAddr().(*net.UDPAddr).IP.String(),
		IntranetPort: listener.LocalAddr().(*net.UDPAddr).Port,
		ExternalIp:   externalUDPAddr.IP.String(),
		ExternalPort: externalUDPAddr.Port,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	m, err := msg.ReadMsg(stream)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	RemoteNetInfo := m.(*models.RemoteNetInfo)
	SendPackToPeer(listener, RemoteNetInfo)
	//TODO：发送认证码用于后续校验
	//TODO:这里控制连接的处理？
	//开始转kcp监听
	kcplis, err := kcp.ServeConn(nil, 10, 3, listener)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	kcplis.SetDeadline(time.Now().Add(time.Second * 5))
	//为了防范风险，只接受一个kcp请求
	//for {
	log.Println("start p2p kcp accpet")
	kcpconn, err := kcplis.AcceptKCP()
	if err != nil {
		if kcplis != nil {
			kcplis.Close()
		}
		if kcpconn != nil {
			kcpconn.Close()
		}
		log.Println(err.Error())
		return nil, err
	}
	nettool.SetYamuxConn(kcpconn)
	log.Println("accpeted")
	log.Println(kcpconn.RemoteAddr())
	//	从从conn中读取p2p另一方发来的认证消息，认证成功之后包装为mux服务端
	err = kcplis.SetDeadline(time.Time{})
	if err != nil {
		kcplis.Close()
		return nil, err
	}
	rawMsg, err := msg.ReadMsgWithTimeOut(kcpconn, time.Second*3)
	if err != nil {
		kcpconn.Close()
		log.Println(err.Error())
		return nil, err
	}
	switch m := rawMsg.(type) {
	//TODO:初步使用ping、pong握手，下一步应该弄成验证校验身份
	case *models.Ping:
		{
			fmt.Printf("P2P握手ping")
			_ = m
			msg.WriteMsg(kcpconn, &models.Pong{})
			config := yamux.DefaultConfig()
			//config.EnableKeepAlive = false
			session, err := yamux.Server(kcpconn, config)
			if err != nil {
				log.Println(err.Error())
			}
			fmt.Printf("Client作为Serverp2p打洞成功！")
			return session, err
		}
	default:
		log.Println("获取到了一个未知的P2P握手消息")
		kcpconn.Close()
		return nil, fmt.Errorf("获取到了一个未知的P2P握手消息")
	}
}
