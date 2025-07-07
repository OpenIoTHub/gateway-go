package services

import (
	"errors"
	"fmt"
	"github.com/OpenIoTHub/gateway-go/v2/utils/qr"
	"github.com/OpenIoTHub/utils/models"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"log"
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
		log.Println(err.Error())
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
func (gm *GatewayCtl) IndexHandler(c *gin.Context) {
	//显示添加的二维码
	htmlContent := `
<!DOCTYPE html>
<html lang="zh">
<head>
    <meta charset="UTF-8">
	<title>OpenIoThub gateway-go - NAT tool for remote control</title>
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
    <div class="tip">使用<a href="https://m.malink.cn/s/RNzqia">云亿连</a>(从应用市场搜索下载或拷贝本链接在移动端打开)扫描上述二维码添加本网关，然后添加主机，主机下面添加端口就可以访问目标端口了！<a href="https://www.bilibili.com/video/BV1Tw9pYJE4B">视频教程🌐</a><a href="https://docs.iothub.cloud/typical/index.html#casaoszimaos">文档🌐</a><a href="https://github.com/OpenIoTHub/gateway-go/v2">开源地址🌐</a></div>
    <div class="tip">Use <a href="https://github.com/OpenIoTHub/OpenIoTHub">OpenIoTHub</a> to scan the following QR code and add a gateway,then add host,add host's port,finally, enjoy remote control.<a href="https://github.com/OpenIoTHub/gateway-go/v2">HomePage🌐</a></div>
</body>
</html>
`
	c.Data(200, "text/html", []byte(htmlContent))
}

// DisplayQrHandler 返回二维码
func (gm *GatewayCtl) DisplayQrHandler(c *gin.Context) {
	var err error
	//显示添加的二维码
	if len(gm.serverSession) == 0 {
		c.Data(200, "text/plain", []byte("no gateway login"))
		return
	}
	gatewayUUID, serverHost, err := gm.GetLoginInfo()
	if err != nil {
		c.Data(200, "text/plain", []byte(err.Error()))
		return
	}

	var qrCode *qrcode.QRCode
	if serverHost == "" || serverHost == qr.STDHost {
		qrCode, err = qr.GetQrById(gatewayUUID)
	} else {
		qrCode, err = qr.GetQrByIdAndHost(gatewayUUID, serverHost)
	}
	if err != nil {
		c.Data(200, "text/plain", []byte(err.Error()))
		return
	}
	c.Header("ContentType", "image/png")
	qrCode.Write(300, c.Writer)
}

// DisplayQrHandler 返回二维码
func (gm *GatewayCtl) GetLoginInfo() (gatewayUUID, serverHost string, err error) {
	for key, sess := range gm.serverSession {
		gatewayUUID = key
		serverHost = sess.tokenModel.Host
	}
	if gatewayUUID == "" && serverHost == "" {
		err = errors.New("Not Logged In")
	}
	return
}
