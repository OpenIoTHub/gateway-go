package gateway

import (
	"errors"
	"fmt"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	nettool "github.com/OpenIoTHub/utils/net"
	"github.com/OpenIoTHub/utils/net/p2p"
	"github.com/libp2p/go-yamux"
	"github.com/xtaci/kcp-go/v5"
	"log"
	"net"
	"time"
)

func MakeP2PSessionAsServer(stream net.Conn, ctrlmMsg *models.ReqNewP2PCtrlAsServer, token *models.TokenClaims) (*yamux.Session, error) {
	if stream != nil {
		defer stream.Close()
	} else {
		return nil, errors.New("stream is nil")
	}
	//监听一个随机端口号，接受P2P方的连接
	externalUDPAddr, listener, err := p2p.GetP2PListener(token)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	p2p.SendPackToPeerByReqNewP2PCtrlAsServer(listener, ctrlmMsg)

	//TODO：发送认证码用于后续校验
	err = msg.WriteMsg(stream, externalUDPAddr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	listener.Close()
	time.Sleep(time.Second)
	//开始转kcp监听
	return kcpListener(listener.LocalAddr().(*net.UDPAddr))
}

//TODO：listener转kcp服务侦听
func kcpListener(laddr *net.UDPAddr) (*yamux.Session, error) {
	//kcplis, err := kcp.ServeConn(nil, 10, 3, listener)
	kcplis, err := kcp.ListenWithOptions(laddr.String(), nil, 10, 3)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
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
		return nil, err
	}
	//配置
	nettool.SetYamuxConn(kcpconn)

	log.Println("accpeted")
	log.Println(kcpconn.RemoteAddr())
	//	从从conn中读取p2p另一方发来的认证消息，认证成功之后包装为mux服务端
	err = kcplis.SetDeadline(time.Time{})
	if err != nil {
		kcplis.Close()
	}
	return kcpConnHdl(kcpconn)
}

func kcpConnHdl(kcpconn net.Conn) (*yamux.Session, error) {
	rawMsg, err := msg.ReadMsgWithTimeOut(kcpconn, time.Second*5)
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
			err = msg.WriteMsg(kcpconn, &models.Pong{})
			if err != nil {
				kcpconn.Close()
				log.Println(err.Error())
				return nil, err
			}
			config := yamux.DefaultConfig()
			//config.EnableKeepAlive = false
			session, err := yamux.Server(kcpconn, config)
			if err != nil {
				log.Println(err.Error())
				return nil, err
			}
			return session, err
		}
	default:
		log.Println("获取到了一个未知的P2P握手消息")
		kcpconn.Close()
		return nil, fmt.Errorf("获取到了一个未知的P2P握手消息")
	}
}
