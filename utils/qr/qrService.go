package qr

import (
	"fmt"
	qrcode "github.com/skip2/go-qrcode"
	"log"
)

const (
	GatewayQRByIdForAddTemplate        = "https://iothub.cloud/a/g?id=%s"
	GatewayQRByIdAndHostForAddTemplate = "https://iothub.cloud/a/g?id=%s&host=%s"
	STDHost                            = "guonei.servers.iothub.cloud"
)

// 通过jwt展示二维码
func DisplayQRCodeById(id, host string) (err error) {
	qrContent := fmt.Sprintf(GatewayQRByIdAndHostForAddTemplate, id, host)
	if host == "" || host == STDHost {
		qrContent = fmt.Sprintf(GatewayQRByIdForAddTemplate, id)
	} else {
		qrContent = fmt.Sprintf(GatewayQRByIdAndHostForAddTemplate, id, host)
	}
	qrCode, err := qrcode.New(qrContent, qrcode.Low)
	if err != nil {
		log.Println(err)
		return
	}
	ascii := qrCode.ToSmallString(false)
	fmt.Println(fmt.Sprintf("Use OpenIoTHub to scan the following QR code and add a gateway(%s)", id))
	fmt.Println(ascii)
	fmt.Println("If the above QR code cannot be scanned, please open the following link and scan the QR code in page:")
	if host == "" || host == STDHost {
		fmt.Println(fmt.Sprintf("https://api.iot-manager.iothub.cloud/v1/displayGatewayQRCodeById?id=%s", id))
	} else {
		fmt.Println(fmt.Sprintf("https://api.iot-manager.iothub.cloud/v1/displayGatewayQRCodeById?id=%s&host=%s", id, host))
	}
	return
}

func GetQrById(id string) (qr *qrcode.QRCode, err error) {
	qrStr := fmt.Sprintf(GatewayQRByIdForAddTemplate, id)
	return qrcode.New(qrStr, qrcode.Low)
}

func GetQrByIdAndHost(id, host string) (qr *qrcode.QRCode, err error) {
	qrStr := fmt.Sprintf(GatewayQRByIdAndHostForAddTemplate, id, host)
	return qrcode.New(qrStr, qrcode.Low)
}
