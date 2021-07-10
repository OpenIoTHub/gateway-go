package services

import (
	"errors"
	"github.com/OpenIoTHub/utils/models"
	"log"
)

var GatewayManager = &GatewayCtl{serverSession: make(map[string]*ServerSession)}

type GatewayCtl struct {
	serverSession map[string]*ServerSession
}

func (gm *GatewayCtl) Loged() bool {
	if len(gm.serverSession) > 0 {
		return true
	}
	return false
}

func (gm *GatewayCtl) AddServer(token string) (err error) {
	tokenModel, err := models.DecodeUnverifiedToken(token)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	if _, ok := gm.serverSession[tokenModel.RunId]; ok {
		log.Println("runId already exist")
		return errors.New("runId already exist")
	}
	serverSession := &ServerSession{
		token:      token,
		tokenModel: tokenModel,
	}
	gm.serverSession[tokenModel.RunId] = serverSession
	return serverSession.start()
}
