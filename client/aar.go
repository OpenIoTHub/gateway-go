package client

import (
	"context"
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
)

type LoginManager struct{}

var loginManager = &LoginManager{}

func Run() {
	s := grpc.NewServer()
	pb.RegisterGatewayLoginManagerServer(s, loginManager)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.GRpcAddr, config.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}
	addr := lis.Addr().(*net.TCPAddr)
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
	_, err = zeroconf.Register(fmt.Sprintf("OpenIoTHubGateway-%s", config.ConfigMode.LastId[:7]), "_openiothub-gateway._tcp", "local.", port,
		[]string{"name=网关",
			"model=com.iotserv.services.gateway",
			fmt.Sprintf("mac=%s", Mac),
			fmt.Sprintf("id=%s", config.ConfigMode.LastId),
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

	//string ConnectionType = 3;
	config.ConfigMode.ConnectionType = in.ConnectionType
	//string LastId = 4;
	config.ConfigMode.LastId = in.LastId

	config.ConfigMode.Server = &models.Srever{
		ServerHost: in.ServerHost,
		TcpPort:    int(in.TcpPort),
		KcpPort:    int(in.KcpPort),
		UdpApiPort: int(in.UdpApiPort),
		KcpApiPort: int(in.KcpApiPort),
		TlsPort:    int(in.TlsPort),
		GrpcPort:   int(in.GrpcPort),
		LoginKey:   in.LoginKey,
	}

	if config.ConfigMode.LastId == "" {
		config.ConfigMode.LastId = uuid.Must(uuid.NewV4()).String()
	}

	GateWayToken, err := models.GetToken(config.ConfigMode, 1, 200000000000)
	if err != nil {
		return &pb.LoginResponse{
			Code:        1,
			Message:     err.Error(),
			LoginStatus: config.Loged,
		}, err
	}
	err = services.GatewayManager.AddServer(GateWayToken)
	if err != nil {
		return &pb.LoginResponse{
			Code:        1,
			Message:     err.Error(),
			LoginStatus: config.Loged,
		}, err
	}
	config.Loged = true
	config.OpenIoTHubToken, err = models.GetToken(config.ConfigMode, 2, 200000000000)
	err = config.WriteConfigFile(config.ConfigMode, config.ConfigFilePath)
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
	OpenIoTHubToken, err := models.GetToken(config.ConfigMode, 2, 200000000000)
	if err != nil {
		return &pb.Token{}, err
	}
	return &pb.Token{Value: OpenIoTHubToken}, nil
}

//rpc GetGateWayToken (Empty) returns (Token) {}
func (lm *LoginManager) GetGateWayToken(ctx context.Context, in *pb.Empty) (*pb.Token, error) {
	return nil, nil
}
