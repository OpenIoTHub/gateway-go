package client

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/OpenIoTHub/gateway-go/component"
	"github.com/OpenIoTHub/gateway-go/config"
	"github.com/OpenIoTHub/gateway-go/services"
	"github.com/OpenIoTHub/gateway-grpc-api/pb-go"
	"github.com/OpenIoTHub/utils/models"
	"github.com/iotdevice/zeroconf"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
)

type LoginManager struct{}

var ConfigMode *models.GatewayConfig
var loginManager = &LoginManager{}

func Run() {
	port, err := strconv.Atoi(config.Setting["gRpcPort"])
	if err != nil {
		log.Println(err)
		return
	}
	//mDNS注册服务
	_, err = zeroconf.Register("OpenIoTHubGateway", "_openiothub-gateway._tcp", "local.", port, []string{}, nil)
	if err != nil {
		log.Println(err)
		return
	}
	//
	s := grpc.NewServer()
	pb.RegisterGatewayLoginManagerServer(s, loginManager)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.Setting["gRpcAddr"], config.Setting["gRpcPort"]))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

//rpc LoginServerByServerInfo (ServerInfo) returns (LoginResponse) {}
func (lm *LoginManager) LoginServerByServerInfo(ctx context.Context, in *pb.ServerInfo) (*pb.LoginResponse, error) {
	var err error
	if config.Loged {
		return &pb.LoginResponse{
			Code:        0,
			Message:     "已经处于登录状态",
			LoginStatus: true,
		}, nil
	}
	//string ServerHost = 1;
	ConfigMode.Server.ServerHost = in.ServerHost
	//string LoginKey = 2;
	ConfigMode.Server.LoginKey = in.LoginKey
	//string ConnectionType = 3;
	ConfigMode.ConnectionType = in.ConnectionType
	//string LastId = 4;
	ConfigMode.LastId = in.LastId
	//int32 TcpPort = 5;
	ConfigMode.Server.TcpPort = int(in.TcpPort)
	//int32 KcpPort = 6;
	ConfigMode.Server.KcpPort = int(in.KcpPort)
	//int32 UdpApiPort = 7;
	ConfigMode.Server.UdpApiPort = int(in.UdpApiPort)
	//int32 KcpApiPort = 8;
	ConfigMode.Server.KcpApiPort = int(in.KcpApiPort)
	//int32 TlsPort = 9;
	ConfigMode.Server.TlsPort = int(in.TlsPort)
	//int32 GrpcPort = 10;
	ConfigMode.Server.GrpcPort = int(in.GrpcPort)

	if ConfigMode.LastId == "" {
		ConfigMode.LastId = uuid.Must(uuid.NewV4()).String()
	}

	GateWayToken, err := models.GetToken(*ConfigMode, 1, 200000000000)
	if err != nil {
		return &pb.LoginResponse{
			Code:        1,
			Message:     err.Error(),
			LoginStatus: config.Loged,
		}, err
	}
	err = services.RunNATManager(ConfigMode.Server.LoginKey, GateWayToken)
	if err != nil {
		return &pb.LoginResponse{
			Code:        1,
			Message:     err.Error(),
			LoginStatus: config.Loged,
		}, err
	}
	config.Loged = true
	config.Setting["OpenIoTHubToken"], err = models.GetToken(*ConfigMode, 2, 200000000000)
	err = config.WriteConfigFile(*ConfigMode, config.Setting["configFilePath"])
	if err != nil {
		log.Println(err.Error())
		return &pb.LoginResponse{
			Code:        1,
			Message:     err.Error(),
			LoginStatus: config.Loged,
		}, err
	}
	return &pb.LoginResponse{
		Code:        0,
		Message:     "登录成功！",
		LoginStatus: config.Loged,
	}, nil
}

//rpc LoginServerByToken (Token) returns (LoginResponse) {}
func (lm *LoginManager) LoginServerByToken(ctx context.Context, in *pb.Token) (*pb.LoginResponse, error) {
	return nil, nil
}

//rpc GetOpenIoTHubToken (Empty) returns (Token) {}
func (lm *LoginManager) GetOpenIoTHubToken(ctx context.Context, in *pb.Empty) (*pb.Token, error) {
	if config.Loged != true || config.Setting["OpenIoTHubToken"] == "" {
		return &pb.Token{}, errors.New("还未登录")
	}
	return &pb.Token{Value: config.Setting["OpenIoTHubToken"]}, nil
}

//rpc GetGateWayToken (Empty) returns (Token) {}
func (lm *LoginManager) GetGateWayToken(ctx context.Context, in *pb.Empty) (*pb.Token, error) {
	return nil, nil
}
