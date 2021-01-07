package services

import (
	"github.com/OpenIoTHub/utils/models"
	"github.com/libp2p/go-yamux"
	"log"
	"sync"
	"time"
)

type ServerSession struct {
	token      string
	tokenModel *models.TokenClaims
	session    *yamux.Session
	heartbeat  *time.Ticker
	quit       chan bool
	sync.Mutex
}

func (ss *ServerSession) stop() {
	ss.quit <- true
}

func (ss *ServerSession) start() (err error) {
	err = ss.LoginServer()
	if err != nil {
		log.Println(err)
		return
	}
	ss.heartbeat = time.NewTicker(time.Second * 20)
	go ss.LoopStream()
	go ss.Task()
	return
}

func (ss *ServerSession) LoginServer() (err error) {
	ss.Lock()
	defer ss.Unlock()
	if ss.session != nil && !ss.session.IsClosed() {
		return
	}
	ss.session, err = LoginServer(ss.token)
	if err != nil {
		log.Println("登录失败：" + err.Error())
		return err
	}
	return
}

func (ss *ServerSession) LoopStream() {
	defer func() {
		if ss.session != nil {
			err := ss.session.Close()
			if err != nil {
				log.Println(err.Error())
			}
		}
	}()
	for {
		if ss.session == nil {
			log.Println("ss.session is nil")
			break
		}
		// Accept a stream
		stream, err := ss.session.AcceptStream()
		if err != nil {
			log.Println("accpStreamErr：" + err.Error())
			break
		}
		log.Println("获取到一个连接需要处理")
		go dlstream(stream, ss.tokenModel, ss.token)
	}
}

func (ss *ServerSession) CheckSessionStatus() {
	if ss.session == nil || (ss.session != nil && ss.session.IsClosed()) {
		err := ss.LoginServer()
		if err != nil {
			log.Println(err)
			return
		}
		ss.LoopStream()
	}
}

func (ss *ServerSession) Task() {
	for {
		select {
		//心跳来了，检测连接的存活状态
		case <-ss.heartbeat.C:
			go ss.CheckSessionStatus()
		case <-ss.quit:
			ss.heartbeat.Stop()
			close(ss.quit)
			if ss.session != nil && !ss.session.IsClosed() {
				log.Println(ss.session.Close())
			}
			return
		}
	}
}
