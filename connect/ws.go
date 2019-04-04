package connect

import (
	"git.iotserv.com/iotserv/utils/io"
	"golang.org/x/net/websocket"
	"net"
)

func JoinWs(stream net.Conn, url, prot, orig string) error {
	ws, err := websocket.Dial(url, prot, orig)
	if err != nil {
		return err
	}
	go io.Join(stream, ws)
	return nil
}

func JoinWss(stream net.Conn, url, prot, orig string) error {
	ws, err := websocket.Dial(url, prot, orig)
	if err != nil {
		return err
	}
	go io.Join(stream, ws)
	return nil
}
