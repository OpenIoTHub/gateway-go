package qr

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/OpenIoTHub/gateway-go/config"
	"github.com/OpenIoTHub/gateway-go/services"
	pb "github.com/OpenIoTHub/openiothub_grpc_api/pb-go/proto/manager"
	qrcode "github.com/skip2/go-qrcode"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net/url"
)

const IoTManagerAddr = "api.iot-manager.iothub.cloud:50051"

// 自动创建jwt并登陆，并展示二维码
func AutoLoginAndDisplayQRCode() (err error) {
	tlsConfig := grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{}))
	conn, err := grpc.NewClient(IoTManagerAddr, tlsConfig)
	if err != nil {
		log.Println("grpc.NewClient:", err)
		return
	}
	defer conn.Close()
	c := pb.NewPublicApiClient(conn)
	md := metadata.Pairs()
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	rst, err := c.GenerateJwtQRCodePair(ctx, &emptypb.Empty{})
	if err != nil {
		log.Println(err)
		return
	}
	err = services.GatewayManager.AddServer(rst.GatewayJwt)
	if err != nil {
		log.Println(err)
		return
	}
	qrs, err := url.ParseRequestURI(rst.QRCodeForMobileAdd)
	if err != nil {
		log.Println(err)
		return
	}
	runId := qrs.Query().Get("id")
	if runId == "" {
		err = errors.New("url id is empty in QRCodeForMobileAdd")
		return
	}
	config.ConfigMode.LoginWithTokenMap[runId] = rst.GatewayJwt
	err = config.WriteConfigFile(config.ConfigMode, config.ConfigFilePath)
	if err != nil {
		log.Println(err)
	}
	return DisplayQRCodeById(runId)
}

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
