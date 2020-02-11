package services

import (
	"fmt"
	"github.com/OpenIoTHub/gateway-go/connect"
	"github.com/OpenIoTHub/utils/crypto"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	//"github.com/OpenIoTHub/utils/io"
	"github.com/jacobsa/go-serial/serial"
	"net"
	"time"

	//"github.com/xtaci/smux"
	"github.com/OpenIoTHub/utils/mux"
)

var lastSalt, lastToken string

func dlstream(stream net.Conn, tokenModel *crypto.TokenClaims) {
	var err error
	defer func() {
		if err == nil || stream == nil {
			return
		}
		err = stream.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	rawMsg, err := msg.ReadMsg(stream)
	if err != nil {
		fmt.Printf(err.Error() + "从stream读取数据错误")
		return
	}
	//fmt.Printf("begin Swc")
	switch m := rawMsg.(type) {
	case *models.ConnectTCP:
		{
			fmt.Printf("tcp")
			err = connect.JoinTCP(stream, m.TargetIP, m.TargetPort)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	case *models.ConnectSTCP:
		{
			fmt.Printf("stcp")
			err = connect.JoinSTCP(stream, m.TargetIP, m.TargetPort)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	case *models.ConnectUDP:
		{
			fmt.Printf("udp")
			err = connect.JoinUDP(stream, m.TargetIP, m.TargetPort)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	case *models.ConnectSerialPort:
		{
			fmt.Printf("sertp")
			err = connect.JoinSerialPort(stream, serial.OpenOptions(*m))
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}

	case *models.ConnectWs:
		{
			fmt.Printf("wstp")
			err = connect.JoinWs(stream, m.TargetUrl, m.Protocol, m.Origin)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}

	case *models.ConnectWss:
		{
			fmt.Printf("wsstp")
			err = connect.JoinWss(stream, m.TargetUrl, m.Protocol, m.Origin)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	case *models.ConnectSSH:
		{
			fmt.Printf("ssh")
			err = connect.JoinSSH(stream, m.TargetIP, m.TargetPort, m.UserName, m.PassWord)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	case *models.NewService:
		{
			fmt.Printf("service")
			err = serviceHdl(stream, m)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	case *models.NewSubSession:
		{
			//:TODO 新创建一个全新的子连接
			fmt.Printf("newSubSession")
			//snappyConn, err := cryptoSnappy.Convert(stream, []byte("BUDIS**$(&CHSKCNNCJSH"))
			//if err != nil {
			//	fmt.Printf(err.Error())
			//	stream.Close()
			//	return
			//}
			config := mux.DefaultConfig()
			//config.EnableKeepAlive = false
			session, err := mux.Server(stream, config)
			if err != nil {
				stream.Close()
				return
			}
			go dlSubSession(session, tokenModel)
		}

	case *models.RequestNewWorkConn:
		{
			fmt.Println("server请求一个新的工作连接")
			stream.Close()
			go newWorkConn(tokenModel)
		}

	case *models.Ping:
		{
			//fmt.Printf("Ping from server")
			err = msg.WriteMsg(stream, &models.Pong{})
			if err != nil {
				fmt.Println(err.Error())
			}
		}

	case *models.ReqNewP2PCtrl:
		{
			fmt.Printf("作为listener方式从洞中获取kcp连接")
			go NewP2PCtrlAsServer(stream, m, tokenModel)
			//lastPing = time.Now()
			//TODO:NETINFO
			//msg.WriteMsg(stream,&models.RemoteNetInfo{
			//	IntranetIp:net.GetIntranetIp(),
			//	IntranetPort:7003,
			//	ExternalIp:net.GetExternalIp(),
			//	ExternalPort:ExternalPort,
			//})
			//stream.Close()
		}
	case *models.ReqNewP2PCtrlAsClient:
		{
			fmt.Printf("作为dial方式从从洞中创建kcp连接")
			go MakeP2PSessionAsClient(stream, tokenModel)
		}
	//	获取检查TCP或者UDP端口状态的请求
	case *models.CheckStatusRequest:
		{
			//fmt.Println("CheckStatusRequest")
			switch m.Type {
			case "tcp", "udp", "tls":
				{
					code, message := connect.CheckTcpUdpTls(m.Type, m.Addr)
					err := msg.WriteMsg(stream, &models.CheckStatusResponse{
						Code:    code,
						Message: message,
					})
					if err != nil {
						fmt.Println(err.Error())
					}
				}
			default:
				err := msg.WriteMsg(stream, &models.CheckStatusResponse{
					Code:    1,
					Message: "type not support",
				})
				if err != nil {
					fmt.Println(err.Error())
				}
			}
			stream.Close()
		}
	default:
		fmt.Printf("type err")
	}
}

func dlsession(session *mux.Session, tokenModel *crypto.TokenClaims) {
	defer func() {
		if session != nil {
			err := session.Close()
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		go func() {
			for {
				err := RunNATManager(lastSalt, lastToken)
				if err != nil {
					fmt.Printf("重新登录失败！原因：%s,5秒钟后重试...\n", err.Error())
					time.Sleep(time.Second * 5)
					continue
				}
				break
			}
		}()
	}()
	for {
		// Accept a stream
		stream, err := session.AcceptStream()
		if err != nil {
			fmt.Println("accpStreamErr：" + err.Error())
			break
		}
		fmt.Println("获取到一个连接需要处理")
		go dlstream(stream, tokenModel)
	}
}

func dlSubSession(session *mux.Session, tokenModel *crypto.TokenClaims) {
	defer func() {
		if session != nil {
			err := session.Close()
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}()
	//session的keepalive,需要配合服务器
	//go func() {
	//	err := io.CheckSession(session)
	//	if err != nil{
	//		fmt.Println(err.Error())
	//		if session != nil{
	//			session.Close()
	//		}
	//	}
	//}()
	for {
		// Accept a stream
		stream, err := session.AcceptStream()
		if err != nil {
			fmt.Printf("accpStream" + err.Error())
			break
		}
		//fmt.Println("Sub Session获取到一个stream处理")
		go dlstream(stream, tokenModel)
	}
	fmt.Printf("exit sub session")
}

func newWorkConn(tokenModel *crypto.TokenClaims) {
	conn, err := LoginWorkConn(tokenModel)
	if err != nil {
		fmt.Println("创建一个到服务端的新的工作连接失败：")
		fmt.Println(err.Error())
		return
	}
	fmt.Println("创建一个到服务端的新的工作连接成功！")
	go dlstream(conn, tokenModel)
}

func RunNATManager(salt, token string) (err error) {
	var session *mux.Session
	var tokenModel *crypto.TokenClaims
	lastSalt, lastToken = salt, token
	session, _, tokenModel, err = Login(salt, token)
	if err != nil {
		//fmt.Println("登录失败：" + err.Error())
		return err
	}
	go dlsession(session, tokenModel)
	return nil
}
