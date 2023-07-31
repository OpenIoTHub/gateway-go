package services

import (
	"github.com/OpenIoTHub/gateway-go/login"
	"github.com/OpenIoTHub/utils/models"
	"github.com/libp2p/go-yamux"
	"log"
	"sync"
	"time"
)

type ServerSession struct {
	token          string
	tokenModel     *models.TokenClaims
	session        *yamux.Session
	heartbeat      *time.Ticker
	quit           chan struct{}
	loginLock      sync.Mutex
	loopStreamLock sync.Mutex
}

func (ss *ServerSession) stop() {
	ss.quit <- struct{}{}
}

func (ss *ServerSession) start() (err error) {
	ss.CheckSessionStatus()
	ss.heartbeat = time.NewTicker(time.Second * 20)
	ss.quit = make(chan struct{})
	go ss.Task()
	return
}

func (ss *ServerSession) LoginServer() (err error) {
	ss.loginLock.Lock()
	defer ss.loginLock.Unlock()
	if ss.session != nil && !ss.session.IsClosed() {
		return
	}
	ss.session, err = login.LoginServer(ss.token)
	if err != nil {
		log.Println("登录失败：" + err.Error())
		return err
	}
	return
}

func (ss *ServerSession) LoopStream() {
	ss.loopStreamLock.Lock()
	defer ss.loopStreamLock.Unlock()
	defer func() {
		if ss.session != nil {
			err := ss.session.Close()
			if err != nil {
				log.Println(err.Error())
			}
		}
	}()
	for {
		if ss.session == nil || (ss.session != nil && ss.session.IsClosed()) {
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
		go handleStream(stream, ss.token)
	}
}

func (ss *ServerSession) CheckSessionStatus() {
	if ss.session == nil || (ss.session != nil && ss.session.IsClosed()) {
		log.Println("开始(重新)连接:", ss.tokenModel.RunId, "@", ss.tokenModel.Host)
		err := ss.LoginServer()
		if err != nil {
			log.Println(err)
			return
		}
		go ss.LoopStream()
	}
}

func (ss *ServerSession) Task() {
	for {
		select {
		//心跳来了，检测连接的存活状态
		case <-ss.heartbeat.C:
			ss.CheckSessionStatus()
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
