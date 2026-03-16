package login

import (
	"github.com/OpenIoTHub/gateway-go/v2/info"
	"github.com/OpenIoTHub/utils/v2/models"
	"github.com/OpenIoTHub/utils/v2/msg"
	"github.com/libp2p/go-yamux"
	"log"
	"net"
	"runtime"
	"strconv"
	"time"
)

func LoginServer(tokenstr string) (*yamux.Session, error) {
	token, err := models.DecodeUnverifiedToken(tokenstr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(token.Host, strconv.Itoa(token.TcpPort)), time.Second*2)
	if err != nil {
		return nil, err
	}
	login := &models.GatewayLogin{
		Token:   tokenstr,
		Os:      runtime.GOOS,
		Arch:    runtime.GOARCH,
		Version: info.Version,
	}
	if err := msg.WriteMsg(conn, login); err != nil {
		conn.Close()
		return nil, err
	}
	session, err := yamux.Server(conn, yamux.DefaultConfig())
	if err != nil {
		conn.Close()
		return nil, err
	}
	log.Println("login OK!")
	return session, nil
}

func LoginWorkConn(tokenStr string) (net.Conn, error) {
	token, err := models.DecodeUnverifiedToken(tokenStr)
	if err != nil {
		return nil, err
	}
	conn, err := net.Dial("tcp", net.JoinHostPort(token.Host, strconv.Itoa(token.TcpPort)))
	if err != nil {
		return nil, err
	}
	loginWorkConn := &models.GatewayWorkConn{
		RunId:   token.RunId,
		Secret:  tokenStr,
		Version: info.Version,
	}
	if err := msg.WriteMsg(conn, loginWorkConn); err != nil {
		conn.Close()
		return nil, err
	}
	return conn, nil
}
