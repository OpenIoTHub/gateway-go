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
	var qrContent string
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
	fmt.Printf("Use OpenIoTHub to scan the following QR code and add a gateway(%s)\n", id)
	fmt.Println(ascii)
	fmt.Println("If the above QR code cannot be scanned, please open the following link and scan the QR code in page:")
	if host == "" || host == STDHost {
		fmt.Printf("https://api.iot-manager.iothub.cloud/v1/displayGatewayQRCodeById?id=%s\n", id)
	} else {
		fmt.Printf("https://api.iot-manager.iothub.cloud/v1/displayGatewayQRCodeById?id=%s&host=%s\n", id, host)
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
