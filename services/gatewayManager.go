package services

import (
	"errors"
	"fmt"
	"github.com/OpenIoTHub/utils/models"
	"log"
)

var GatewayManager = &GatewayCtl{serverSession: make(map[string]*ServerSession)}

type GatewayCtl struct {
	serverSession map[string]*ServerSession
}

func (gm *GatewayCtl) Loged() bool {
	return len(gm.serverSession) > 0
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

func (gm *GatewayCtl) DelServer(runid string) (err error) {
	if _, ok := gm.serverSession[runid]; ok {
		log.Println("找到了runid的serverSession")
		gm.serverSession[runid].stop()
		delete(gm.serverSession, runid)
		//TODO 同时删除配置文件的相关配置
		return
	}
	return errors.New(fmt.Sprintf("gateway uuid:%s not found", runid))
}
