package handle

import (
	connect "github.com/OpenIoTHub/gateway-go/v2/netservice/services/connect/conn"
	"github.com/OpenIoTHub/gateway-go/v2/netservice/services/connect/service"
	"github.com/OpenIoTHub/gateway-go/v2/netservice/services/login"
	"github.com/OpenIoTHub/utils/v2/models"
	"github.com/OpenIoTHub/utils/v2/msg"
	"github.com/OpenIoTHub/utils/v2/net/p2p/gateway"
	"log"
	"net"

	"github.com/libp2p/go-yamux"
)

func HandleStream(stream net.Conn, tokenStr string) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic HandleStream: %+v", err)
		}
	}()
	var tokenModel *models.TokenClaims
	if tokenStr != "" {
		var err error
		tokenModel, err = models.DecodeUnverifiedToken(tokenStr)
		if err != nil {
			log.Printf("解析token失败: %v", err)
			// 继续执行，某些操作可能不需要token
		}
	}
	rawMsg, err := msg.ReadMsg(stream)
	if err != nil {
		log.Printf("从stream读取数据错误: %v", err)
		stream.Close()
		return
	}
	switch m := rawMsg.(type) {
	case *models.ConnectTCP:
		log.Printf("处理TCP连接: %s:%d", m.TargetIP, m.TargetPort)
		if err := connect.JoinTCP(stream, m.TargetIP, m.TargetPort); err != nil {
			log.Printf("TCP连接失败: %v", err)
			return
		}
	case *models.ConnectSTCP:
		log.Printf("处理STCP连接: %s:%d", m.TargetIP, m.TargetPort)
		if err := connect.JoinSTCP(stream, m.TargetIP, m.TargetPort); err != nil {
			log.Printf("STCP连接失败: %v", err)
			return
		}
	case *models.ConnectUDP:
		log.Printf("处理UDP连接: %s:%d", m.TargetIP, m.TargetPort)
		if err := connect.JoinUDP(stream, m.TargetIP, m.TargetPort); err != nil {
			log.Printf("UDP连接失败: %v", err)
			return
		}
	case *models.ConnectSerialPort:
		log.Printf("处理串口连接")
		if err := connect.JoinSerialPort(stream, m); err != nil {
			log.Printf("串口连接失败: %v", err)
			return
		}
	case *models.ConnectWs:
		log.Printf("处理WebSocket连接: %s", m.TargetUrl)
		if err := connect.JoinWs(stream, m.TargetUrl, m.Protocol, m.Origin); err != nil {
			log.Printf("WebSocket连接失败: %v", err)
			return
		}
	case *models.ConnectWss:
		log.Printf("处理WebSocket Secure连接: %s", m.TargetUrl)
		if err := connect.JoinWss(stream, m.TargetUrl, m.Protocol, m.Origin); err != nil {
			log.Printf("WebSocket Secure连接失败: %v", err)
			return
		}
	case *models.ConnectSSH:
		log.Printf("处理SSH连接: %s:%d", m.TargetIP, m.TargetPort)
		if err := connect.JoinSSH(stream, m.TargetIP, m.TargetPort, m.UserName, m.PassWord); err != nil {
			log.Printf("SSH连接失败: %v", err)
			return
		}
	case *models.NewService:
		if err := service.ServiceHdl(stream, m); err != nil {
			log.Printf("处理新服务失败: %v", err)
			return
		}
	case *models.NewSubSession:
		log.Printf("创建新的子会话")
		config := yamux.DefaultConfig()
		session, err := yamux.Server(stream, config)
		if err != nil {
			log.Printf("创建yamux会话失败: %v", err)
			stream.Close()
			return
		}
		go HandleSession(session, tokenStr)

	case *models.RequestNewWorkConn:
		log.Println("服务器请求一个新的工作连接")
		stream.Close()
		go newWorkConn(tokenStr)

	case *models.Ping:
		if err := msg.WriteMsg(stream, &models.Pong{}); err != nil {
			log.Printf("发送Pong失败: %v", err)
		}
		//TODO 防止未关闭的连接，取决于请求方是否关闭
		//stream.Close()

	case *models.ReqNewP2PCtrlAsServer:
		log.Printf("作为listener方式从洞中获取kcp连接")
		if tokenModel == nil {
			log.Println("tokenModel为空，无法创建P2P会话")
			stream.Close()
			return
		}
		go func() {
			session, listener, err := gateway.MakeP2PSessionAsServer(stream, m, tokenModel)
			if err != nil {
				if listener != nil {
					listener.Close()
				}
				log.Printf("创建P2P服务器会话失败: %v", err)
				return
			}
			defer func() {
				if listener != nil {
					listener.Close()
				}
			}()
			HandleSession(session, tokenStr)
		}()
	case *models.ReqNewP2PCtrlAsClient:
		log.Printf("作为dial方式从洞中创建kcp连接")
		if tokenModel == nil {
			log.Println("tokenModel为空，无法创建P2P会话")
			stream.Close()
			return
		}
		go func() {
			session, listener, err := gateway.MakeP2PSessionAsClient(stream, m, tokenModel)
			if err != nil {
				if listener != nil {
					listener.Close()
				}
				log.Printf("创建P2P客户端会话失败: %v", err)
				return
			}
			defer func() {
				if listener != nil {
					listener.Close()
				}
			}()
			HandleSession(session, tokenStr)
		}()
	//	获取检查TCP或者UDP端口状态的请求
	case *models.CheckStatusRequest:
		var response *models.CheckStatusResponse
		switch m.Type {
		case "tcp", "udp", "tls":
			code, message := service.CheckTcpUdpTls(m.Type, m.Addr)
			response = &models.CheckStatusResponse{
				Code:    code,
				Message: message,
			}
		default:
			response = &models.CheckStatusResponse{
				Code:    1,
				Message: "type not support",
			}
		}
		if err := msg.WriteMsg(stream, response); err != nil {
			log.Printf("发送检查状态响应失败: %v", err)
		}
		stream.Close()
	//由于用户在服务器账户删掉了这个网关，所有网关删掉服务器登录以供新用户绑定
	case *models.DeleteGatewayJwt:
		//	TODO: 实现删除网关JWT的逻辑
		stream.Close()
	default:
		log.Printf("未知的消息类型: %T", rawMsg)
		stream.Close()
	}
}

func HandleSession(session *yamux.Session, tokenStr string) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic HandleSession: %+v", err)
		}
		if session != nil {
			if err := session.Close(); err != nil {
				log.Printf("关闭session失败: %v", err)
			}
		}
	}()
	for {
		if session == nil {
			return
		}
		stream, err := session.AcceptStream()
		if err != nil {
			log.Printf("从session接受流失败: %v", err)
			if stream != nil {
				stream.Close()
			}
			break
		}
		go HandleStream(stream, tokenStr)
	}
}

// 新创建的工作连接
func newWorkConn(tokenStr string) {
	if tokenStr == "" {
		log.Println("token为空，无法创建工作连接")
		return
	}
	conn, err := login.LoginWorkConn(tokenStr)
	if err != nil {
		log.Printf("创建到服务端的工作连接失败: %v", err)
		if conn != nil {
			conn.Close()
		}
		return
	}
	log.Println("创建到服务端的工作连接成功")
	go HandleStream(conn, tokenStr)
}
