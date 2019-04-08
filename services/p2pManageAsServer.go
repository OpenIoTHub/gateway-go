package services

import (
	"fmt"
	"git.iotserv.com/iotserv/utils/crypto"
	"git.iotserv.com/iotserv/utils/models"
	"git.iotserv.com/iotserv/utils/msg"
	"git.iotserv.com/iotserv/utils/mux"
	"git.iotserv.com/iotserv/utils/net"
	"github.com/xtaci/kcp-go"
	"net"
)

//var ExternalPort int
//var SendPackReqChan = make(chan *models.SendUdpPackReq,10)

func NewP2PCtrlAsServer(stream net.Conn, ctrlmMsg *models.ReqNewP2PCtrl, token *crypto.TokenClaims) {
	//监听一个随机端口号，接受P2P方的连接
	localIps, localPort, externalIp, externalPort, listener, err := nettool.GetP2PListener(token)
	if err != nil {
		fmt.Println(err)
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
func kcpListener(listener *net.UDPConn, token *crypto.TokenClaims) {
	kcplis, err := kcp.ServeConn(nil, 10, 3, listener)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	//为了防范风险，只接受一个kcp请求
	//for {
	fmt.Println("start p2p kcp accpet")
	kcpconn, err := kcplis.AcceptKCP()
	//配置
	kcpconn.SetStreamMode(true)
	kcpconn.SetWriteDelay(false)
	kcpconn.SetNoDelay(0, 100, 1, 1)
	kcpconn.SetWindowSize(128, 256)
	kcpconn.SetMtu(1350)
	kcpconn.SetACKNoDelay(true)

	fmt.Println("accpeted")
	fmt.Println(kcpconn.RemoteAddr())
	if err != nil {
		fmt.Println(err.Error())
	}
	//b:=make([]byte,1024)
	//n,err:=conn.Read(b)
	//fmt.Println(string(b[0:n]))
	//lis.Close()
	//	从从conn中读取p2p另一方发来的认证消息，认证成功之后包装为mux服务端
	go kcpConnHdl(kcpconn, token)
	//}
}

func kcpConnHdl(kcpconn net.Conn, token *crypto.TokenClaims) {
	rawMsg, err := msg.ReadMsg(kcpconn)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	switch m := rawMsg.(type) {
	//TODO:初步使用ping、pong握手，下一步应该弄成验证校验身份
	case *models.Ping:
		{
			fmt.Printf("P2P握手ping")
			_ = m
			msg.WriteMsg(kcpconn, &models.Pong{})
			config := mux.DefaultConfig()
			session, err := mux.Server(kcpconn, config)
			if err != nil {
				fmt.Println(err.Error())
			}
			go dlSubSession(session, token)
			fmt.Printf("Client作为Serverp2p打洞成功！")
		}
	default:
		fmt.Println("获取到了一个未知的P2P握手消息")
	}
}
