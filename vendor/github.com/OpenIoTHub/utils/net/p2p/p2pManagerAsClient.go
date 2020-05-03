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

//作为客户端主动去连接内网client的方式创建穿透连接
func MakeP2PSessionAsClient(stream net.Conn, token *models.TokenClaims) (*yamux.Session, error) {
	if stream != nil {
		defer stream.Close()
	}
	ExternalUDPAddr, listener, err := GetP2PListener(token)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	msgsd := &models.ReqNewP2PCtrl{
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
	case *models.RemoteNetInfo:
		{
			fmt.Printf("remote net info")
			//TODO:认证；同内网直连；抽象出公共函数？
			kcpconn, err := kcp.NewConn(fmt.Sprintf("%s:%d", m.ExternalIp, m.ExternalPort), nil, 10, 3, listener)
			nettool.SetYamuxConn(kcpconn)
			if err != nil {
				fmt.Printf(err.Error())
				return nil, err
			}
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
					fmt.Printf("get pong from p2p kcpconn")
					_ = m
					//TODO:认证
					config := yamux.DefaultConfig()
					//config.EnableKeepAlive = false
					p2pSubSession, err := yamux.Server(kcpconn, config)
					if err != nil {
						if p2pSubSession != nil {
							p2pSubSession.Close()
						}
						fmt.Printf("create sub session err:" + err.Error())
						return nil, err
					}
					//return p2pSubSession
					return p2pSubSession, err
				}
			default:
				fmt.Printf("type err")
			}
		}
	default:
		fmt.Printf("type err")
	}
	return nil, err
}
