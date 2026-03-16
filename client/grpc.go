package client

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/OpenIoTHub/gateway-go/v2/config"
	"github.com/OpenIoTHub/gateway-go/v2/services"
	pb "github.com/OpenIoTHub/openiothub_grpc_api/pb-go/proto/gateway"
	"github.com/OpenIoTHub/utils/v2/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

var GRpcPort = 55443

type LoginManager struct {
	pb.UnimplementedGatewayLoginManagerServer
}

func startGRPC() {
	s := grpc.NewServer()
	pb.RegisterGatewayLoginManagerServer(s, &LoginManager{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.GRpcAddr, GRpcPort))
	if err != nil {
		log.Printf("gRPC 监听失败: %v", err)
		return
	}
	log.Printf("Grpc 监听端口: %d\n", GRpcPort)
	reflection.Register(s)
	go registerGatewayMDNS(GRpcPort)
	if err := s.Serve(lis); err != nil {
		log.Printf("gRPC 服务失败: %v", err)
	}
}

func (lm *LoginManager) CheckGatewayLoginStatus(ctx context.Context, in *emptypb.Empty) (*pb.LoginResponse, error) {
	return &pb.LoginResponse{
		Code:        0,
		Message:     "网关登录状态",
		LoginStatus: services.GatewayManager.Loged(),
	}, nil
}

func (lm *LoginManager) LoginServerByToken(ctx context.Context, in *pb.Token) (*pb.LoginResponse, error) {
	if services.GatewayManager.Loged() && !IsLibrary {
		return &pb.LoginResponse{
			Code:        1,
			Message:     "网关已经登录服务器",
			LoginStatus: services.GatewayManager.Loged(),
		}, nil
	}
	tokenModel, err := models.DecodeUnverifiedToken(in.Value)
	if err != nil {
		log.Println(err.Error())
		return &pb.LoginResponse{
			Code:        1,
			Message:     "token错误",
			LoginStatus: services.GatewayManager.Loged(),
		}, err
	}
	if err := services.GatewayManager.AddServer(in.Value); err != nil {
		return &pb.LoginResponse{
			Code:        1,
			Message:     err.Error(),
			LoginStatus: services.GatewayManager.Loged(),
		}, nil
	}
	config.ConfigMode.LoginWithTokenMap[tokenModel.RunId] = in.Value
	if err := config.WriteConfigFile(config.ConfigMode, config.ConfigFilePath); err != nil {
		log.Printf("保存配置文件失败: %v", err)
	}
	return &pb.LoginResponse{
		Code:        0,
		Message:     "登录成功！",
		LoginStatus: services.GatewayManager.Loged(),
	}, nil
}

func (lm *LoginManager) testEmbeddedByValue() {}
