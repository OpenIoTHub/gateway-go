package qr

import (
	"fmt"
	qrcode "github.com/skip2/go-qrcode"
	"log"
)

const (
	IoTManagerAddr          = "api.iot-manager.iothub.cloud:50051"
	GatewayQRForAddTemplate = "https://iothub.cloud/a/g?id=%s"
)

// 通过jwt展示二维码
func DisplayQRCodeById(id string) (err error) {
	qrContent := fmt.Sprintf("https://iothub.cloud/a/g?id=%s", id)
	qrCode, err := qrcode.New(qrContent, qrcode.Low)
	if err != nil {
		log.Println(err)
		return
	}
	ascii := qrCode.ToSmallString(false)
	fmt.Println(fmt.Sprintf("Use OpenIoTHub to scan the following QR code and add a gateway(%s)", id))
	fmt.Println(ascii)
	fmt.Println("If the above QR code cannot be scanned, please open the following link and scan the QR code in page:")
	fmt.Println(fmt.Sprintf("https://api.iot-manager.iothub.cloud/v1/displayGatewayQRCodeById?id=%s", id))
	return
}

func GetQrById(id string) (qr *qrcode.QRCode, err error) {
	qrStr := fmt.Sprintf(GatewayQRForAddTemplate, id)
	return qrcode.New(qrStr, qrcode.Low)
}
