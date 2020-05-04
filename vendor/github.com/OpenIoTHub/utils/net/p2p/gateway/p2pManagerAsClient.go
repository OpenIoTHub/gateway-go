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

//作为客户端主动去连接内网client的方式创建穿透连接
func MakeP2PSessionAsClient(stream net.Conn, ctrlmMsg *models.ReqNewP2PCtrlAsClient, token *models.TokenClaims) (*yamux.Session, error) {
	if stream != nil {
		defer stream.Close()
	} else {
		return nil, errors.New("stream is nil")
	}
	ExternalUDPAddr, listener, err := p2p.GetP2PListener(token)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	msgsd := &models.RemoteNetInfo{
		IntranetIp:   listener.LocalAddr().(*net.UDPAddr).IP.String(),
		IntranetPort: listener.LocalAddr().(*net.UDPAddr).Port,
		ExternalIp:   ExternalUDPAddr.IP.String(),
		ExternalPort: ExternalUDPAddr.Port,
	}
	err = msg.WriteMsg(stream, msgsd)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rawMsg, err := msg.ReadMsg(stream)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	switch m := rawMsg.(type) {
	case *models.OK:
		{
			_ = m
			log.Printf("remote net info")
			//TODO:认证；同内网直连；抽象出公共函数？
			kcpconn, err := kcp.NewConn(fmt.Sprintf("%s:%d", ctrlmMsg.ExternalIp, ctrlmMsg.ExternalPort), nil, 10, 3, listener)
			if err != nil {
				log.Printf(err.Error())
				return nil, err
			}
			//设置
			nettool.SetYamuxConn(kcpconn)
			time.Sleep(time.Second)
			err = msg.WriteMsg(kcpconn, &models.Ping{})
			if err != nil {
				kcpconn.Close()
				log.Println(err)
				return nil, err
			}

			rawMsg, err := msg.ReadMsgWithTimeOut(kcpconn, time.Second*3)
			if err != nil {
				kcpconn.Close()
				log.Println(err)
				return nil, err
			}
			switch m := rawMsg.(type) {
			case *models.Pong:
				{
					log.Printf("get pong from p2p kcpconn")
					_ = m
					//TODO:认证
					config := yamux.DefaultConfig()
					//config.EnableKeepAlive = false
					p2pSubSession, err := yamux.Server(kcpconn, config)
					if err != nil {
						if p2pSubSession != nil {
							p2pSubSession.Close()
						}
						log.Printf("create sub session err:" + err.Error())
						return nil, err
					}
					//return p2pSubSession
					return p2pSubSession, err
				}
			default:
				log.Printf("type err")
			}
		}
	default:
		log.Printf("type err")
	}
	return nil, errors.New("gateway p2pManagerAsClient 失败")
}
