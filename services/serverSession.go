package services

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/OpenIoTHub/gateway-go/v2/netservice/handle"
	"github.com/OpenIoTHub/gateway-go/v2/netservice/services/login"
	"github.com/OpenIoTHub/utils/v2/models"
	"github.com/libp2p/go-yamux"
)

type ServerSession struct {
	token      string
	tokenModel *models.TokenClaims

	session   *yamux.Session
	heartbeat *time.Ticker
	quit      chan struct{}

	checkSessionStatusLock sync.Mutex
	loginLock              sync.Mutex
	loopStreamLock         sync.Mutex
}

func (ss *ServerSession) stop() {
	select {
	case ss.quit <- struct{}{}:
	default:
	}
}

func (ss *ServerSession) start() error {
	ss.quit = make(chan struct{}, 1)
	ss.heartbeat = time.NewTicker(20 * time.Second)
	ss.checkSessionStatus()
	go ss.task()
	return nil
}

func (ss *ServerSession) loginServer() error {
	ss.loginLock.Lock()
	defer ss.loginLock.Unlock()
	if ss.session != nil && !ss.session.IsClosed() {
		return nil
	}
	session, err := login.LoginServer(ss.token)
	if err != nil {
		return fmt.Errorf("登录失败: %w", err)
	}
	ss.session = session
	return nil
}

func (ss *ServerSession) loopStream() {
	ss.loopStreamLock.Lock()
	defer ss.loopStreamLock.Unlock()
	for {
		if ss.session == nil || ss.session.IsClosed() {
			log.Println("session is nil or closed")
			break
		}
		stream, err := ss.session.AcceptStream()
		if err != nil {
			log.Printf("接受流失败: %v", err)
			if ss.session != nil {
				ss.session.Close()
			}
			break
		}
		go handle.HandleStream(stream, ss.token)
	}
}

func (ss *ServerSession) checkSessionStatus() {
	ss.checkSessionStatusLock.Lock()
	defer ss.checkSessionStatusLock.Unlock()
	if ss.session == nil || ss.session.IsClosed() {
		log.Printf("开始(重新)连接: %s @ %s", ss.tokenModel.RunId, ss.tokenModel.Host)
		if err := ss.loginServer(); err != nil {
			log.Printf("检查会话状态时登录失败: %v", err)
			return
		}
		go ss.loopStream()
	}
}

func (ss *ServerSession) task() {
	defer func() {
		ss.heartbeat.Stop()
		if ss.session != nil && !ss.session.IsClosed() {
			if err := ss.session.Close(); err != nil {
				log.Printf("关闭session失败: %v", err)
			}
		}
	}()
	for {
		select {
		case <-ss.heartbeat.C:
			ss.checkSessionStatus()
		case <-ss.quit:
			return
		}
	}
}
