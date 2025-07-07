package handle

import (
	connect "github.com/OpenIoTHub/gateway-go/v2/netservice/services/connect/conn"
	"github.com/OpenIoTHub/gateway-go/v2/netservice/services/connect/service"
	"github.com/OpenIoTHub/gateway-go/v2/netservice/services/login"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	"github.com/OpenIoTHub/utils/net/p2p/gateway"
	"log"

	"net"
	//"github.com/xtaci/smux"
	"github.com/libp2p/go-yamux"
)

func HandleStream(stream net.Conn, tokenStr string) {
	var err error
	tokenModel, err := models.DecodeUnverifiedToken(tokenStr)
	if err != nil {
		log.Println(err.Error())
		return
	}
	rawMsg, err := msg.ReadMsg(stream)
	if err != nil {
		log.Println(err.Error() + "从stream读取数据错误")
		stream.Close()
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
			err = connect.JoinSerialPort(stream, m)
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
			//log.Printf("case *models.NewService")
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
			go HandleSession(session, tokenStr)
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
			//TODO 防止未关闭的连接，取决于请求方是否关闭
			//stream.Close()
		}

	case *models.ReqNewP2PCtrlAsServer:
		{
			log.Printf("作为listener方式从洞中获取kcp连接")
			go func() {
				session, listener, err := gateway.MakeP2PSessionAsServer(stream, m, tokenModel)
				if err != nil {
					if listener != nil {
						listener.Close()
					}
					log.Println("gateway.MakeP2PSessionAsServer:", err)
					return
				}
				HandleSession(session, tokenStr)
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
					if listener != nil {
						listener.Close()
					}
					log.Println("gateway.MakeP2PSessionAsClient:", err)
					return
				}
				HandleSession(session, tokenStr)
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
					code, message := service.CheckTcpUdpTls(m.Type, m.Addr)
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
			//TODO 是否关闭
			stream.Close()
		}
	//由于用户在服务器账户删掉了这个网关，所有网关删掉服务器登录以供新用户绑定
	case *models.DeleteGatewayJwt:
		{
			//	log.Println("删除配置:", tokenModel.RunId)
			//	GatewayManager.DelServer(tokenModel.RunId)
			//	delete(ConfigMode.LoginWithTokenMap, tokenModel.RunId)
			//	err = WriteConfigFile(ConfigMode, ConfigFilePath)
			//	if err != nil {
			//		log.Println(err)
			//		err = msg.WriteMsg(stream, &models.Error{
			//			Code:    1,
			//			Message: err.Error(),
			//		})
			//		if err != nil {
			//			log.Println(err.Error())
			//		}
			//		return
			//	}
			//	err = msg.WriteMsg(stream, &models.OK{})
			//	if err != nil {
			//		log.Println(err.Error())
			//	}
			stream.Close()
		}
	default:
		log.Printf("type err")
		stream.Close()
	}
}

func HandleSession(session *yamux.Session, tokenStr string) {
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
			log.Println("accept stream form session got err：" + err.Error())
			if stream != nil {
				stream.Close()
			}
			break
		}
		//log.Println("获取到一个连接需要处理")
		go HandleStream(stream, tokenStr)
	}
}

// 新创建的工作连接
func newWorkConn(tokenStr string) {
	conn, err := login.LoginWorkConn(tokenStr)
	if err != nil {
		log.Println("创建一个到服务端的新的工作连接失败：")
		log.Println(err.Error())
		if conn != nil {
			conn.Close()
		}
		return
	}
	log.Println("创建一个到服务端的新的工作连接成功！")
	go HandleStream(conn, tokenStr)
}
