package services

import (
	"errors"
	"fmt"
	"github.com/OpenIoTHub/gateway-go/utils/qr"
	"github.com/OpenIoTHub/utils/models"
	"log"
	"net/http"
)

var GatewayManager = &GatewayCtl{serverSession: make(map[string]*ServerSession)}

type GatewayCtl struct {
	serverSession map[string]*ServerSession
}

func (gm *GatewayCtl) Loged() bool {
	return len(gm.serverSession) > 0
}

// AddServer 添加网关实例，登录一个id
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

// DelServer 删除网关实例，删除一个id
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

// IndexHandler http服务首页
func (gm *GatewayCtl) IndexHandler(w http.ResponseWriter, r *http.Request) {
	//显示添加的二维码
	w.Header().Set("Content-Type", "text/html")
	htmlContent := `
<!DOCTYPE html>
<html lang="zh">
<head>
    <meta charset="UTF-8">
    <style>
        body {
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            margin: 0;
        }
        img {
            max-width: 100%;
            height: auto;
            margin-bottom: 20px;
        }
        .tip {
            color: green;
            text-align: center;
            font-size: 1.2em;
        }
    </style>
</head>
<body>
    <img src="/DisplayQrHandler" alt="扫码添加二维码">
    <div class="tip">使用<a href="https://m.malink.cn/s/RNzqia">云亿连</a>(从应用市场下载)扫描上述二维码添加本网关</div>
    <div class="tip">Use <a href="https://github.com/OpenIoTHub/OpenIoTHub">OpenIoTHub</a> to scan the following QR code and add a gateway</div>
</body>
</html>
`
	fmt.Fprintf(w, htmlContent)
}

// DisplayQrHandler 返回二维码
func (gm *GatewayCtl) DisplayQrHandler(w http.ResponseWriter, r *http.Request) {
	//显示添加的二维码
	if len(gm.serverSession) == 0 {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "no gateway login")
		return
	}
	gatewayUUID := ""
	for key, _ := range gm.serverSession {
		gatewayUUID = key
	}

	qr, err := qr.GetQrById(gatewayUUID)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, err.Error())
		return
	}
	w.Header().Set("ContentType", "image/png")
	qr.Write(300, w)
}
