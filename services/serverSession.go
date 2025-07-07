package services

import (
	"github.com/OpenIoTHub/gateway-go/v2/netservice/handle"
	"github.com/OpenIoTHub/gateway-go/v2/netservice/services/login"
	"github.com/OpenIoTHub/utils/models"
	"github.com/libp2p/go-yamux"
	"log"
	"sync"
	"time"
)

type ServerSession struct {
	//基础信息
	token      string
	tokenModel *models.TokenClaims
	//内部存储
	session   *yamux.Session
	heartbeat *time.Ticker
	quit      chan struct{}

	checkSessionStatusLock sync.Mutex
	loginLock              sync.Mutex
	loopStreamLock         sync.Mutex
}

func (ss *ServerSession) stop() {
	ss.quit <- struct{}{}
}

func (ss *ServerSession) start() (err error) {
	//防止多次调用
	ss.checkSessionStatus()
	ss.heartbeat = time.NewTicker(time.Second * 20)
	ss.quit = make(chan struct{})
	go ss.task()
	return
}

func (ss *ServerSession) loginServer() (err error) {
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

func (ss *ServerSession) loopStream() {
	ss.loopStreamLock.Lock()
	defer ss.loopStreamLock.Unlock()
	//防止影响新创建的会话，不关闭会话
	//defer func() {
	//	if ss.session != nil {
	//		err := ss.session.Close()
	//		if err != nil {
	//			log.Println(err.Error())
	//		}
	//	}
	//}()
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
		go handle.HandleStream(stream, ss.token)
	}
}

func (ss *ServerSession) checkSessionStatus() {
	ss.checkSessionStatusLock.Lock()
	defer ss.checkSessionStatusLock.Unlock()
	if ss.session == nil || (ss.session != nil && ss.session.IsClosed()) {
		log.Println("开始(重新)连接:", ss.tokenModel.RunId, "@", ss.tokenModel.Host)
		err := ss.loginServer()
		if err != nil {
			log.Println(err)
			return
		}
		go ss.loopStream()
	}
}

func (ss *ServerSession) task() {
	for {
		select {
		//心跳来了，检测连接的存活状态
		case <-ss.heartbeat.C:
			ss.checkSessionStatus()
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
