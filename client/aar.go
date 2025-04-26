package client

import (
	"context"
	"fmt"
	"github.com/OpenIoTHub/gateway-go/config"
	"github.com/OpenIoTHub/gateway-go/info"
	"github.com/OpenIoTHub/gateway-go/services"
	"github.com/OpenIoTHub/openiothub_grpc_api/pb-go/proto/gateway"
	"github.com/OpenIoTHub/utils/models"
	"github.com/grandcat/zeroconf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"net/http"
)

var IsLibrary = true

type LoginManager struct {
	pb.UnimplementedGatewayLoginManagerServer
}

var loginManager = new(LoginManager)

func Run() {
	go start()
}

func start() {
	//启动http服务
	go func() {
		http.HandleFunc("/", services.GatewayManager.IndexHandler)
		http.HandleFunc("/DisplayQrHandler", services.GatewayManager.DisplayQrHandler)
		//不同的系统打印不一样
		log.Printf("Http 监听端口: %d\n", config.ConfigMode.HttpServicePort)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", config.ConfigMode.HttpServicePort), nil); err != nil {
			log.Printf("Failed to start server: %s\n", err)
		}
	}()
	//启动grpc服务
	s := grpc.NewServer()
	pb.RegisterGatewayLoginManagerServer(s, loginManager)
	//port := services.GrpcPort
	//if runtime.GOOS == "android" {
	//	port = 55443
	//}
	port := 55443
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.GRpcAddr, port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}
	//addr := lis.Addr().(*net.TCPAddr)
	log.Printf("Grpc 监听端口:%d\n", port)
	reflection.Register(s)
	go regMDNS(port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func regMDNS(port int) {
	var Mac = "mac"
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println(err)
	} else if len(interfaces) > 0 {
		Mac = interfaces[0].HardwareAddr.String()
	}
	gatewayUUID, serverHost, err := services.GatewayManager.GetLoginInfo()
	//qrStr, err := qr.GetQrByIdAndHost(gatewayUUID, serverHost)
	//mDNS注册服务
	_, err = zeroconf.Register(fmt.Sprintf("OpenIoTHubGateway-%s", config.ConfigMode.GatewayUUID), "_openiothub-gateway._tcp", "local.", port,
		[]string{"name=网关",
			"model=com.iotserv.services.gateway",
			fmt.Sprintf("mac=%s", Mac),
			fmt.Sprintf("id=%s", config.ConfigMode.GatewayUUID),
			//提供网关添加信息
			fmt.Sprintf("RunId=%s", gatewayUUID),
			fmt.Sprintf("ServerHost=%s", serverHost),
			"author=Farry",
			"email=newfarry@126.com",
			"home-page=https://github.com/OpenIoTHub",
			"firmware-respository=https://github.com/OpenIoTHub/gateway-go",
			fmt.Sprintf("firmware-version=%s", info.Version)}, nil)
	if err != nil {
		log.Println(err)
		return
	}
}

// rpc CheckGatewayLoginStatus (Empty) returns (LoginResponse) {}
func (lm *LoginManager) CheckGatewayLoginStatus(ctx context.Context, in *emptypb.Empty) (*pb.LoginResponse, error) {
	return &pb.LoginResponse{
		Code:        0,
		Message:     "网关登录状态",
		LoginStatus: services.GatewayManager.Loged(),
	}, nil
}

// rpc LoginServerByServerInfo (ServerInfo) returns (LoginResponse) {}
func (lm *LoginManager) LoginServerByToken(ctx context.Context, in *pb.Token) (*pb.LoginResponse, error) {
	//如果已经登录则阻止登录
	if services.GatewayManager.Loged() && !IsLibrary {
		return &pb.LoginResponse{
			Code:        1,
			Message:     "网关已经登录服务器",
			LoginStatus: services.GatewayManager.Loged(),
		}, nil
	}
	tokenModel, err := models.DecodeUnverifiedToken(in.Value)
	if err != nil {
		log.Printf(err.Error())
		return &pb.LoginResponse{
			Code:        1,
			Message:     "token错误",
			LoginStatus: services.GatewayManager.Loged(),
		}, err
	}
	//使用token登录
	err = services.GatewayManager.AddServer(in.Value)
	if err != nil {
		return &pb.LoginResponse{
			Code:        1,
			Message:     err.Error(),
			LoginStatus: services.GatewayManager.Loged(),
		}, nil
	}
	config.ConfigMode.LoginWithTokenMap[tokenModel.RunId] = in.Value
	err = config.WriteConfigFile(config.ConfigMode, config.ConfigFilePath)
	if err != nil {
		log.Println(err)
	}
	//标记为已经登录并返回结果
	return &pb.LoginResponse{
		Code:        0,
		Message:     "登录成功！",
		LoginStatus: services.GatewayManager.Loged(),
	}, nil
}

func (lm *LoginManager) testEmbeddedByValue() {
}
