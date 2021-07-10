package client

import (
	"context"
	"fmt"
	_ "github.com/OpenIoTHub/gateway-go/component"
	"github.com/OpenIoTHub/gateway-go/config"
	"github.com/OpenIoTHub/gateway-go/services"
	"github.com/OpenIoTHub/gateway-grpc-api/pb-go"
	"github.com/grandcat/zeroconf"
	"google.golang.org/grpc"
	"log"
	"net"
)

type LoginManager struct {
	*pb.UnimplementedGatewayLoginManagerServer
}

var loginManager = new(LoginManager)

func Run() {
	s := grpc.NewServer()
	pb.RegisterGatewayLoginManagerServer(s, loginManager)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.GRpcAddr, config.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}
	addr := lis.Addr().(*net.TCPAddr)
	fmt.Printf("Grpc 监听端口:%d\n", addr.Port)
	go regMDNS(addr.Port)
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
	//mDNS注册服务
	_, err = zeroconf.Register(fmt.Sprintf("OpenIoTHubGateway-%s", config.ConfigMode.GatewayUUID), "_openiothub-gateway._tcp", "local.", port,
		[]string{"name=网关",
			"model=com.iotserv.services.gateway",
			fmt.Sprintf("mac=%s", Mac),
			fmt.Sprintf("id=%s", config.ConfigMode.GatewayUUID),
			"author=Farry",
			"email=newfarry@126.com",
			"home-page=https://github.com/OpenIoTHub",
			"firmware-respository=https://github.com/OpenIoTHub/gateway-go",
			fmt.Sprintf("firmware-version=%s", services.Version)}, nil)
	if err != nil {
		log.Println(err)
		return
	}
}

//rpc CheckGatewayLoginStatus (Empty) returns (LoginResponse) {}
func (lm *LoginManager) CheckGatewayLoginStatus(ctx context.Context, in *pb.Empty) (*pb.LoginResponse, error) {
	return &pb.LoginResponse{
		Code:        0,
		Message:     "网关登录状态",
		LoginStatus: services.GatewayManager.Loged(),
	}, nil
}

//rpc LoginServerByServerInfo (ServerInfo) returns (LoginResponse) {}
func (lm *LoginManager) LoginServerByToken(ctx context.Context, in *pb.Token) (*pb.LoginResponse, error) {
	//如果已经登录则阻止登录
	if services.GatewayManager.Loged() {
		return &pb.LoginResponse{
			Code:        1,
			Message:     "网关已经登录服务器",
			LoginStatus: services.GatewayManager.Loged(),
		}, nil
	}
	//使用token登录
	err := services.GatewayManager.AddServer(in.Value)
	if err != nil {
		return &pb.LoginResponse{
			Code:        1,
			Message:     err.Error(),
			LoginStatus: services.GatewayManager.Loged(),
		}, nil
	}
	config.ConfigMode.LoginWithTokenList = append(config.ConfigMode.LoginWithTokenList, in.Value)
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
