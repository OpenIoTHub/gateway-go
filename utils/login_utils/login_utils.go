package login_utils

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/OpenIoTHub/gateway-go/config"
	"github.com/OpenIoTHub/gateway-go/services"
	"github.com/OpenIoTHub/gateway-go/utils/qr"
	pb "github.com/OpenIoTHub/openiothub_grpc_api/pb-go/proto/manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net/url"
)

const (
	IoTManagerAddr = "api.iot-manager.iothub.cloud:50051"
)

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
	return qr.DisplayQRCodeById(runId)
}
