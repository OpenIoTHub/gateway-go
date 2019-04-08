package services

import (
	"fmt"
	"git.iotserv.com/iotserv/client/connect"
	"git.iotserv.com/iotserv/client/connect/serial"
	"git.iotserv.com/iotserv/utils/crypto"
	"git.iotserv.com/iotserv/utils/models"
	"git.iotserv.com/iotserv/utils/msg"
	"net"

	//"github.com/xtaci/smux"
	"git.iotserv.com/iotserv/utils/mux"
	"time"
)

var mytoken string
var lastPing time.Time
var try = true
var sub = false

func dlstream(stream net.Conn, tokenModel *crypto.TokenClaims) {
	rawMsg, err := msg.ReadMsg(stream)
	if err != nil {
		fmt.Printf(err.Error() + "从stream读取数据错误")
		stream.Close()
		return
	}
	//fmt.Printf("begin Swc")
	switch m := rawMsg.(type) {
	case *models.ConnectTCP:
		{
			fmt.Printf("tcp")
			err := connect.JoinTCP(stream, m.TargetIP, m.TargetPort)
			if err != nil {
				stream.Close()
			}
		}
	case *models.ConnectSTCP:
		{
			fmt.Printf("stcp")
			err := connect.JoinSTCP(stream, m.TargetIP, m.TargetPort)
			if err != nil {
				stream.Close()
			}
		}
	case *models.ConnectUDP:
		{
			fmt.Printf("udp")
			err := connect.JoinUDP(stream, m.TargetIP, m.TargetPort)
			if err != nil {
				stream.Close()
			}
		}
	case *models.ConnectSerialPort:
		{
			fmt.Printf("sertp")
			err := serial.JoinSerialPort(stream, m.TargetPort, m.Baud)
			if err != nil {
				stream.Close()
			}
		}

	case *models.ConnectWs:
		{
			fmt.Printf("wstp")
			err := connect.JoinWs(stream, m.TargetUrl, m.Protocol, m.Origin)
			if err != nil {
				stream.Close()
			}
		}

	case *models.ConnectWss:
		{
			fmt.Printf("wsstp")
			err := connect.JoinWss(stream, m.TargetUrl, m.Protocol, m.Origin)
			if err != nil {
				stream.Close()
			}
		}
	case *models.ConnectSSH:
		{
			fmt.Printf("ssh")
			err := connect.JoinSSH(stream, m.TargetIP, m.TargetPort, m.UserName, m.PassWord)
			if err != nil {
				stream.Close()
			}
		}
	case *models.NewService:
		{
			fmt.Printf("service")
			err := serviceHdl(stream, m)
			if err != nil {
				stream.Close()
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
				return
			}
			sub = true
			go dlSubSession(session, tokenModel)
		}

	case *models.RequestNewWorkConn:
		{
			fmt.Println("server请求一个新的工作连接")
			go newWorkConn(tokenModel)
		}

	case *models.Ping:
		{
			defer func() {
				err := stream.Close()
				if err != nil {
					fmt.Println(err.Error())
				}
			}()
			fmt.Printf("Ping from server")
			lastPing = time.Now()
			err := msg.WriteMsg(stream, &models.Pong{})
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
			defer func() {
				if stream != nil {
					err := stream.Close()
					if err != nil {
						fmt.Println(err.Error())
					}
				}
			}()
			fmt.Println("CheckStatusRequest")
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
		}
	default:
		fmt.Printf("type err")
	}
}

func dlsession(session *mux.Session, tokenModel *crypto.TokenClaims, salt string) {
	relogin := func() {
		for {
			retry, err := RunNATManager(salt, mytoken)
			if err != nil {
				if retry == false { //停止重试登陆
					fmt.Printf("token 过期或者token错误，停止尝试登陆")
					break
				}
				time.Sleep(time.Second * 60)
			} else {
				break
			}
		}
	}
	for {
		// Accept a stream
		stream, err := session.AcceptStream()
		if err != nil {
			fmt.Printf("accpStream" + err.Error())
			fmt.Println(time.Now())
			try = true
			break
		}
		fmt.Println("Session获取到一个stream处理")
		go dlstream(stream, tokenModel)
	}
	if try != false {
		go relogin()
		fmt.Printf("重新登陆")
	} else {
		fmt.Printf("可能是token校验失败，停止尝试登陆")
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
	for {
		// Accept a stream
		stream, err := session.AcceptStream()
		if err != nil {
			fmt.Printf("accpStream" + err.Error())
			fmt.Println(time.Now())
			break
		}
		fmt.Println("Sub Session获取到一个stream处理")
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

func RunNATManager(salt, token string) (bool, error) {
	var session *mux.Session
	var retry bool
	var err error
	var tokenModel *crypto.TokenClaims
	mytoken = token
	session, retry, tokenModel, err = Login(salt, token)
	if err != nil {
		fmt.Printf("login" + err.Error())
		return retry, err
	}
	go dlsession(session, tokenModel, salt)
	return retry, nil
}
