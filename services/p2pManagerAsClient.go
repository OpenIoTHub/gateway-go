package services

import (
	"fmt"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	"github.com/OpenIoTHub/utils/mux"
	"github.com/OpenIoTHub/utils/net"
	"github.com/xtaci/kcp-go/v5"
	"log"
	"net"
	"time"
)

//作为客户端主动去连接内网client的方式创建穿透连接
func MakeP2PSessionAsClient(stream net.Conn, token *models.TokenClaims) {
	if stream != nil {
		defer stream.Close()
	}
	ExternalUDPAddr, listener, err := nettool.GetP2PListener(token)
	if err != nil {
		log.Println(err.Error())
		return
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
		return
	}
	rawMsg, err := msg.ReadMsg(stream)
	if err != nil {
		log.Println(err)
		return
	}
	switch m := rawMsg.(type) {
	case *models.RemoteNetInfo:
		{
			fmt.Printf("remote net info")
			//TODO:认证；同内网直连；抽象出公共函数？
			kcpconn, err := kcp.NewConn(fmt.Sprintf("%s:%d", m.ExternalIp, m.ExternalPort), nil, 10, 3, listener)
			//设置
			kcpconn.SetStreamMode(true)
			kcpconn.SetWriteDelay(false)
			kcpconn.SetNoDelay(0, 100, 1, 1)
			kcpconn.SetWindowSize(128, 256)
			kcpconn.SetMtu(1350)
			kcpconn.SetACKNoDelay(true)
			if err != nil {
				fmt.Printf(err.Error())
				return
			}
			err = msg.WriteMsg(kcpconn, &models.Ping{})
			if err != nil {
				kcpconn.Close()
				log.Println(err)
				return
			}

			rawMsg, err := msg.ReadMsgWithTimeOut(kcpconn, time.Second*3)
			if err != nil {
				kcpconn.Close()
				log.Println(err)
				return
			}
			switch m := rawMsg.(type) {
			case *models.Pong:
				{
					fmt.Printf("get pong from p2p kcpconn")
					_ = m
					//TODO:认证
					config := mux.DefaultConfig()
					//config.EnableKeepAlive = false
					p2pSubSession, err := mux.Server(kcpconn, config)
					if err != nil {
						if p2pSubSession != nil {
							p2pSubSession.Close()
						}
						fmt.Printf("create sub session err:" + err.Error())
						return
					}
					//return p2pSubSession
					go dlSubSession(p2pSubSession, token)
				}
			default:
				fmt.Printf("type err")
			}
		}
	default:
		fmt.Printf("type err")
	}
}
