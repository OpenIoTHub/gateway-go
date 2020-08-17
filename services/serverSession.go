package services

import (
	"github.com/OpenIoTHub/utils/models"
	"github.com/libp2p/go-yamux"
	"log"
	"sync"
	"time"
)

type ServerSession struct {
	token         string
	tokenModel    *models.TokenClaims
	session       *yamux.Session
	heartbeat     *time.Ticker
	quitHeartbeat chan bool
	sync.Mutex
}

func (ss *ServerSession) start() (err error) {
	err = ss.LoginServer()
	if err != nil {
		log.Println(err)
		return
	}
	ss.heartbeat = time.NewTicker(time.Second * 20)
	go ss.LoopStream()
	return
}

func (ss *ServerSession) LoginServer() error {
	return nil
}

func (ss *ServerSession) LoopStream() {

}

func (ss *ServerSession) CheckSessionStatus() {

}

func (ss *ServerSession) Task() {
	for {
		select {
		//心跳来了，检测连接的存活状态
		case <-ss.heartbeat.C:
			go ss.CheckSessionStatus()
		case <-ss.quitHeartbeat:
			ss.heartbeat.Stop()
			close(ss.quitHeartbeat)
			return
		}
	}
}
