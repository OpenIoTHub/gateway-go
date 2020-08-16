package services

import (
	"github.com/OpenIoTHub/utils/models"
	"github.com/libp2p/go-yamux"
	"log"
	"sync"
)

type ServerSession struct {
	token      string
	tokenModel *models.TokenClaims
	session    *yamux.Session
	sync.Mutex
}

func (ss *ServerSession) start() (err error) {
	err = ss.LoginServer()
	if err != nil {
		log.Println(err)
		return
	}
	go ss.LoopStream()
	return
}

func (ss *ServerSession) LoginServer() error {
	return nil
}

func (ss *ServerSession) LoopStream() {

}
