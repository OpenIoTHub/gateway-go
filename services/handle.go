package services

import (
	"github.com/OpenIoTHub/gateway-go/netservice/connect"
	"github.com/OpenIoTHub/gateway-go/netservice/login"
	"github.com/OpenIoTHub/gateway-go/netservice/service"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	"github.com/OpenIoTHub/utils/net/p2p/gateway"
	"log"

	//"github.com/OpenIoTHub/utils/io"
	"github.com/jacobsa/go-serial/serial"
	"net"
	//"github.com/xtaci/smux"
	"github.com/libp2p/go-yamux"
)

func handleStream(stream net.Conn, tokenStr string) {
	var err error
	tokenModel, err := models.DecodeUnverifiedToken(tokenStr)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	defer func() {
		if err == nil || stream == nil {
			return
		}
		err = stream.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}()
	rawMsg, err := msg.ReadMsg(stream)
	if err != nil {
		log.Printf(err.Error() + "从stream读取数据错误")
		return
	}
	//log.Printf("begin Swc")
	switch m := rawMsg.(type) {
	case *models.ConnectTCP:
		{
			log.Printf("tcp")
			err = connect.JoinTCP(stream, m.TargetIP, m.TargetPort)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	case *models.ConnectSTCP:
		{
			log.Printf("stcp")
			err = connect.JoinSTCP(stream, m.TargetIP, m.TargetPort)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	case *models.ConnectUDP:
		{
			log.Printf("udp")
			err = connect.JoinUDP(stream, m.TargetIP, m.TargetPort)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	case *models.ConnectSerialPort:
		{
			log.Printf("sertp")
			err = connect.JoinSerialPort(stream, serial.OpenOptions(*m))
			if err != nil {
				log.Println(err.Error())
				return
			}
		}

	case *models.ConnectWs:
		{
			log.Printf("wstp")
			err = connect.JoinWs(stream, m.TargetUrl, m.Protocol, m.Origin)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}

	case *models.ConnectWss:
		{
			log.Printf("wsstp")
			err = connect.JoinWss(stream, m.TargetUrl, m.Protocol, m.Origin)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	case *models.ConnectSSH:
		{
			log.Printf("ssh")
			err = connect.JoinSSH(stream, m.TargetIP, m.TargetPort, m.UserName, m.PassWord)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	case *models.NewService:
		{
			log.Printf("case *models.NewService")
			err = service.ServiceHdl(stream, m)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	case *models.NewSubSession:
		{
			//:TODO 新创建一个全新的子连接
			log.Printf("newSubSession")
			//snappyConn, err := modelsSnappy.Convert(stream, []byte("BUDIS**$(&CHSKCNNCJSH"))
			//if err != nil {
			//	log.Printf(err.Error())
			//	stream.Close()
			//	return
			//}
			config := yamux.DefaultConfig()
			//config.EnableKeepAlive = false
			session, err := yamux.Server(stream, config)
			if err != nil {
				stream.Close()
				return
			}
			go handleSession(session, tokenStr)
		}

	case *models.RequestNewWorkConn:
		{
			log.Println("server请求一个新的工作连接")
			stream.Close()
			go newWorkConn(tokenStr)
		}

	case *models.Ping:
		{
			//log.Printf("Ping from server")
			err = msg.WriteMsg(stream, &models.Pong{})
			if err != nil {
				log.Println(err.Error())
			}
		}

	case *models.ReqNewP2PCtrlAsServer:
		{
			log.Printf("作为listener方式从洞中获取kcp连接")
			go func() {
				session, listener, err := gateway.MakeP2PSessionAsServer(stream, m, tokenModel)
				if err != nil {
					log.Println("gateway.MakeP2PSessionAsServer:", err)
					return
				}
				handleSession(session, tokenStr)
				if listener != nil {
					listener.Close()
				}
			}()

		}
	case *models.ReqNewP2PCtrlAsClient:
		{
			log.Printf("作为dial方式从从洞中创建kcp连接")
			go func() {
				session, listener, err := gateway.MakeP2PSessionAsClient(stream, m, tokenModel)
				if err != nil {
					log.Println("gateway.MakeP2PSessionAsClient:", err)
					return
				}
				handleSession(session, tokenStr)
				if listener != nil {
					listener.Close()
				}
			}()
		}
	//	获取检查TCP或者UDP端口状态的请求
	case *models.CheckStatusRequest:
		{
			//log.Println("CheckStatusRequest")
			switch m.Type {
			case "tcp", "udp", "tls":
				{
					code, message := connect.CheckTcpUdpTls(m.Type, m.Addr)
					err := msg.WriteMsg(stream, &models.CheckStatusResponse{
						Code:    code,
						Message: message,
					})
					if err != nil {
						log.Println(err.Error())
					}
				}
			default:
				err := msg.WriteMsg(stream, &models.CheckStatusResponse{
					Code:    1,
					Message: "type not support",
				})
				if err != nil {
					log.Println(err.Error())
				}
			}
			stream.Close()
		}
	//由于用户在服务器账户删掉了这个网关，所有网关删掉服务器登录以供新用户绑定
	case *models.DeleteGatewayJwt:
		{
			log.Println("删除配置:", tokenModel.RunId)
			GatewayManager.DelServer(tokenModel.RunId)
			delete(ConfigMode.LoginWithTokenMap, tokenModel.RunId)
			err = WriteConfigFile(ConfigMode, ConfigFilePath)
			if err != nil {
				log.Println(err)
				err = msg.WriteMsg(stream, &models.Error{
					Code:    1,
					Message: err.Error(),
				})
				if err != nil {
					log.Println(err.Error())
				}
				return
			}
			err = msg.WriteMsg(stream, &models.OK{})
			if err != nil {
				log.Println(err.Error())
			}
		}
	default:
		log.Printf("type err")
	}
}

func handleSession(session *yamux.Session, tokenStr string) {
	defer func() {
		if session != nil {
			err := session.Close()
			if err != nil {
				log.Println(err.Error())
			}
		}
	}()
	for {
		// Accept a stream
		stream, err := session.AcceptStream()
		if err != nil {
			log.Println("accpStreamErr：" + err.Error())
			break
		}
		//log.Println("获取到一个连接需要处理")
		go handleStream(stream, tokenStr)
	}
}

// 新创建的工作连接
func newWorkConn(tokenStr string) {
	conn, err := login.LoginWorkConn(tokenStr)
	if err != nil {
		log.Println("创建一个到服务端的新的工作连接失败：")
		log.Println(err.Error())
		return
	}
	log.Println("创建一个到服务端的新的工作连接成功！")
	go handleStream(conn, tokenStr)
}
