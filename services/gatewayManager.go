package services

import (
	"errors"
	"fmt"
	"sync"

	"github.com/OpenIoTHub/utils/v2/models"
)

var GatewayManager = &GatewayCtl{serverSession: make(map[string]*ServerSession)}

type GatewayCtl struct {
	mu            sync.RWMutex
	serverSession map[string]*ServerSession
}

func (gm *GatewayCtl) Loged() bool {
	gm.mu.RLock()
	defer gm.mu.RUnlock()
	return len(gm.serverSession) > 0
}

func (gm *GatewayCtl) AddServer(token string) error {
	tokenModel, err := models.DecodeUnverifiedToken(token)
	if err != nil {
		return fmt.Errorf("decode token failed: %w", err)
	}
	gm.mu.Lock()
	if _, ok := gm.serverSession[tokenModel.RunId]; ok {
		gm.mu.Unlock()
		return fmt.Errorf("runId %s already exists", tokenModel.RunId)
	}
	ss := &ServerSession{
		token:      token,
		tokenModel: tokenModel,
	}
	gm.serverSession[tokenModel.RunId] = ss
	gm.mu.Unlock()
	return ss.start()
}

func (gm *GatewayCtl) DelServer(runid string) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	session, ok := gm.serverSession[runid]
	if !ok {
		return fmt.Errorf("gateway uuid: %s not found", runid)
	}
	session.stop()
	delete(gm.serverSession, runid)
	return nil
}

func (gm *GatewayCtl) GetLoginInfo() (gatewayUUID, serverHost string, err error) {
	gm.mu.RLock()
	defer gm.mu.RUnlock()
	if len(gm.serverSession) == 0 {
		return "", "", errors.New("not logged in")
	}
	for key, sess := range gm.serverSession {
		return key, sess.tokenModel.Host, nil
	}
	return "", "", errors.New("no active session found")
}
