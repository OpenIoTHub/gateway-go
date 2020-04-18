package services

import (
	"fmt"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	"github.com/OpenIoTHub/utils/mux"
	"net"
	"runtime"
	"strconv"
)

var ServerIp string

func Login(salt, tokenstr string) (*mux.Session, bool, *models.TokenClaims, error) { //bool retry? false :dont retry
	token, err := models.DecodeToken(salt, tokenstr)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, false, &models.TokenClaims{}, err
	}
	ServerIp = token.Host
	//KCP方式
	//conn, err := kcp.DialWithOptions(fmt.Sprintf("%s:%d", ServerIp, token.KcpPort), nil, 10, 3)
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
		return nil, true, token, err
	}
	login := &models.GatewayLogin{
		Token: tokenstr,
		Os:    runtime.GOOS,
		Arch:  runtime.GOARCH,
	}

	err = msg.WriteMsg(conn, login)
	if err != nil {
		conn.Close()
		return nil, true, token, err
	}
	config := mux.DefaultConfig()
	//config.EnableKeepAlive = false
	session, err := mux.Server(conn, config)
	if err != nil {
		conn.Close()
		return nil, false, token, err
	}
	fmt.Printf("login OK!")
	return session, false, token, nil
}

func LoginWorkConn(token *models.TokenClaims) (net.Conn, error) {
	ServerIp = token.Host
	//KCP方式
	//conn, err := kcp.DialWithOptions(fmt.Sprintf("%s:%d", ServerIp, token.KcpPort), nil, 10, 3)
	//conn.SetStreamMode(true)
	//conn.SetWriteDelay(false)
	//conn.SetNoDelay(0, 40, 2, 1)
	//conn.SetWindowSize(1024, 1024)
	//conn.SetMtu(1472)
	//conn.SetACKNoDelay(true)
	//Tls
	//conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", ServerIp, token.TlsPort), &tls.Config{InsecureSkipVerify: true})
	//TCP
	conn, err := net.Dial("tcp", net.JoinHostPort(ServerIp, strconv.Itoa(token.TcpPort)))
	if err != nil {
		return nil, err
	}
	loginWorkConn := &models.GatewayWorkConn{
		RunId:  token.RunId,
		Secret: "",
	}

	err = msg.WriteMsg(conn, loginWorkConn)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
