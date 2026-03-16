package services

import (
	"github.com/OpenIoTHub/gateway-go/v2/utils/qr"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

func (gm *GatewayCtl) IndexHandler(c *gin.Context) {
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
    <div class="tip">使用<a href="https://m.malink.cn/s/RNzqia">云亿连</a>(从应用市场搜索下载或拷贝本链接在移动端打开)扫描上述二维码添加本网关，然后添加主机，主机下面添加端口就可以访问目标端口了！<a href="https://www.bilibili.com/video/BV1Tw9pYJE4B">视频教程🌐</a><a href="https://docs.iothub.cloud/typical/index.html#casaoszimaos">文档🌐</a><a href="https://github.com/OpenIoTHub/gateway-go">开源地址🌐</a></div>
    <div class="tip">Use <a href="https://github.com/OpenIoTHub/OpenIoTHub">OpenIoTHub</a> to scan the above QR code and add a gateway,then add host,add host's port,finally, enjoy remote control.<a href="https://github.com/OpenIoTHub/gateway-go">HomePage🌐</a></div>
</body>
</html>
`
	c.Data(200, "text/html", []byte(htmlContent))
}

func (gm *GatewayCtl) DisplayQrHandler(c *gin.Context) {
	if !gm.Loged() {
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
