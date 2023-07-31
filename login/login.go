package login

import (
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	"github.com/libp2p/go-yamux"
	"log"
	"net"
	"runtime"
	"strconv"
	"time"
)

var (
	Version = "dev"
)

func LoginServer(tokenstr string) (*yamux.Session, error) { //bool retry? false :dont retry
	token, err := models.DecodeUnverifiedToken(tokenstr)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	//KCP方式
	//conn, err := kcp.DialWithOptions(fmt.Sprintf("%s:%d", token.Host, token.KcpPort), nil, 10, 3)
	//conn.SetStreamMode(true)
	//conn.SetWriteDelay(false)
	//conn.SetNoDelay(0, 40, 2, 1)
	//conn.SetWindowSize(1024, 1024)
	//conn.SetMtu(1472)
	//conn.SetACKNoDelay(true)
	//Tls
	//conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", token.Host, token.TlsPort), &tls.Config{InsecureSkipVerify: true})
	//TCP
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(token.Host, strconv.Itoa(token.TcpPort)), time.Second*2)
	if err != nil {
		return nil, err
	}
	login := &models.GatewayLogin{
		Token:   tokenstr,
		Os:      runtime.GOOS,
		Arch:    runtime.GOARCH,
		Version: Version,
	}

	err = msg.WriteMsg(conn, login)
	if err != nil {
		conn.Close()
		return nil, err
	}
	config := yamux.DefaultConfig()
	//config.EnableKeepAlive = false
	session, err := yamux.Server(conn, config)
	if err != nil {
		conn.Close()
		return nil, err
	}
	log.Printf("login OK!")
	return session, nil
}

func LoginWorkConn(tokenStr string) (net.Conn, error) {
	token, err := models.DecodeUnverifiedToken(tokenStr)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	//KCP方式
	//conn, err := kcp.DialWithOptions(fmt.Sprintf("%s:%d", token.Host, token.KcpPort), nil, 10, 3)
	//conn.SetStreamMode(true)
	//conn.SetWriteDelay(false)
	//conn.SetNoDelay(0, 40, 2, 1)
	//conn.SetWindowSize(1024, 1024)
	//conn.SetMtu(1472)
	//conn.SetACKNoDelay(true)
	//Tls
	//conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", token.Host, token.TlsPort), &tls.Config{InsecureSkipVerify: true})
	//TCP
	conn, err := net.Dial("tcp", net.JoinHostPort(token.Host, strconv.Itoa(token.TcpPort)))
	if err != nil {
		return nil, err
	}
	loginWorkConn := &models.GatewayWorkConn{
		RunId:   token.RunId,
		Secret:  tokenStr,
		Version: Version,
	}

	err = msg.WriteMsg(conn, loginWorkConn)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
