package services

import (
	"fmt"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	"github.com/OpenIoTHub/utils/mux"
	"github.com/OpenIoTHub/utils/net"
	"github.com/xtaci/kcp-go"
	"log"
	"net"
	"time"
)

//var ExternalPort int
//var SendPackReqChan = make(chan *models.SendUdpPackReq,10)

func NewP2PCtrlAsServer(stream net.Conn, ctrlmMsg *models.ReqNewP2PCtrl, token *models.TokenClaims) {
	//监听一个随机端口号，接受P2P方的连接
	localIps, localPort, externalIp, externalPort, listener, err := nettool.GetP2PListener(token)
	if err != nil {
		log.Println(err)
		return
	}
	nettool.SendPackToPeer(listener, ctrlmMsg)
	//开始转kcp监听
	go kcpListener(listener, token)
	//TODO：发送认证码用于后续校验
	msg.WriteMsg(stream, &models.RemoteNetInfo{
		IntranetIp:   localIps,
		IntranetPort: localPort,
		ExternalIp:   externalIp,
		ExternalPort: externalPort,
	})
	//TODO:这里控制连接的处理？
	stream.Close()
}

//TODO：listener转kcp服务侦听
func kcpListener(listener *net.UDPConn, token *models.TokenClaims) {
	kcplis, err := kcp.ServeConn(nil, 10, 3, listener)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	kcplis.SetDeadline(time.Now().Add(time.Second * 5))
	//为了防范风险，只接受一个kcp请求
	//for {
	log.Println("start p2p kcp accpet")
	kcpconn, err := kcplis.AcceptKCP()
	if err != nil {
		kcplis.Close()
		if kcpconn != nil {
			kcpconn.Close()
		}
		log.Println(err.Error())
		return
	}
	//配置
	//kcpconn.SetDeadline(time.Now().Add(time.Second * 5))
	kcpconn.SetStreamMode(true)
	kcpconn.SetWriteDelay(false)
	kcpconn.SetNoDelay(0, 100, 1, 1)
	kcpconn.SetWindowSize(128, 256)
	kcpconn.SetMtu(1350)
	kcpconn.SetACKNoDelay(true)

	log.Println("accpeted")
	log.Println(kcpconn.RemoteAddr())
	//b:=make([]byte,1024)
	//n,err:=conn.Read(b)
	//log.Println(string(b[0:n]))
	//lis.Close()
	//	从从conn中读取p2p另一方发来的认证消息，认证成功之后包装为mux服务端
	err = kcplis.SetDeadline(time.Time{})
	if err != nil {
		kcplis.Close()
	}
	err = kcpConnHdl(kcpconn, token)
	if err != nil {
		kcplis.Close()
	}
	//}
}

func kcpConnHdl(kcpconn net.Conn, token *models.TokenClaims) error {
	rawMsg, err := msg.ReadMsgWithTimeOut(kcpconn, time.Second*3)
	if err != nil {
		kcpconn.Close()
		log.Println(err.Error())
		return err
	}
	switch m := rawMsg.(type) {
	//TODO:初步使用ping、pong握手，下一步应该弄成验证校验身份
	case *models.Ping:
		{
			fmt.Printf("P2P握手ping")
			_ = m
			msg.WriteMsg(kcpconn, &models.Pong{})
			config := mux.DefaultConfig()
			//config.EnableKeepAlive = false
			session, err := mux.Server(kcpconn, config)
			if err != nil {
				log.Println(err.Error())
			}
			go dlSubSession(session, token)
			fmt.Printf("Client作为Serverp2p打洞成功！")
			return nil
		}
	default:
		log.Println("获取到了一个未知的P2P握手消息")
		kcpconn.Close()
		return fmt.Errorf("获取到了一个未知的P2P握手消息")
	}
}
