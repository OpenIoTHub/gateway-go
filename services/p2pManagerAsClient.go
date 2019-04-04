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
	"strconv"
	"strings"
)

type connectedUDPConn struct{ *net.UDPConn }

func (c *connectedUDPConn) WriteTo(b []byte, addr net.Addr) (int, error) { return c.Write(b) }

//作为客户端主动去连接内网client的方式创建穿透连接
func MakeP2PSessionAsClient(stream net.Conn, token *crypto.TokenClaims) {
	//stream, err := session.OpenStream()
	//if err != nil {
	//	fmt.Printf("get session" + err.Error())
	//	return nil
	//}
	localAddr, ip, port, err := nettool.GetDialIpPort(token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	localPort, err := strconv.Atoi(strings.Split(localAddr.String(), ":")[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	msgsd := &models.ReqNewP2PCtrl{
		IntranetIp:   nettool.GetIntranetIp(),
		IntranetPort: localPort,
		ExternalIp:   ip,
		ExternalPort: port,
	}
	err = msg.WriteMsg(stream, msgsd)
	if err != nil {
		fmt.Println(err)
		return
	}
	rawMsg, err := msg.ReadMsg(stream)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch m := rawMsg.(type) {
	case *models.RemoteNetInfo:
		{
			fmt.Printf("remote net info")
			//TODO:认证；同内网直连；抽象出公共函数？
			//kcpconn, err := kcp.DialWithOptions(fmt.Sprintf("%s:%d", m.ExternalIp, m.ExternalPort), nil, 10, 3)
			raddr := fmt.Sprintf("%s:%d", m.ExternalIp, m.ExternalPort)
			udpaddr, err := net.ResolveUDPAddr("udp", raddr)
			if err != nil {
				return
			}
			laddr, err := net.ResolveUDPAddr("udp", localAddr.String())
			if err != nil {
				return
			}
			udpconn, err := net.DialUDP("udp", laddr, udpaddr)
			if err != nil {
				return
			}

			kcpconn, err := kcp.NewConn(raddr, nil, 10, 3, &connectedUDPConn{udpconn})
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
				fmt.Println(err)
				return
			}

			rawMsg, err := msg.ReadMsg(kcpconn)
			if err != nil {
				fmt.Println(err)
				return
			}
			switch m := rawMsg.(type) {
			case *models.Pong:
				{
					fmt.Printf("get pong from p2p kcpconn")
					_ = m
					//TODO:认证
					p2pSubSession, err := mux.Server(kcpconn, nil)
					if err != nil {
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
