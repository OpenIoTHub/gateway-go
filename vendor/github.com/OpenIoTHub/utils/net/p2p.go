package nettool

import (
	"fmt"
	"github.com/OpenIoTHub/utils/crypto"
	"github.com/OpenIoTHub/utils/models"
	"net"
	"strconv"
	"strings"
	"time"
)

//获取一个listener的外部地址和端口
func GetExternalIpPort(listener *net.UDPConn, token *crypto.TokenClaims) (ip string, port int, err error) {
	udpaddr, err := net.ResolveUDPAddr("udp", token.Host+":"+strconv.Itoa(token.P2PApiPort))
	//udpaddr, err := net.ResolveUDPAddr("udp", "tencent-shanghai-v1.host.nat-cloud.com:34321")
	if err != nil {
		fmt.Printf("%s", err.Error())
		return "", 0, err
	}

	err = listener.SetDeadline(time.Now().Add(time.Duration(3 * time.Second)))
	if err != nil {
		fmt.Printf("%s", err.Error())
		return "", 0, err
	}

	listener.WriteToUDP([]byte("getIpPort"), udpaddr)

	fmt.Println("发送到服务器确定成功！等待确定外网ip和port")
	data := make([]byte, 256)
	n, _, err := listener.ReadFromUDP(data)
	fmt.Println("获取api的UDP包成功，开始解析自己listener出口地址和端口")
	if err != nil {
		fmt.Printf("获取listener的出口出错: %s", err.Error())
		return "", 0, err
	}
	ipPort := string(data[:n])
	ip = strings.Split(ipPort, ":")[0]
	port, err = strconv.Atoi(strings.Split(ipPort, ":")[1])
	if err != nil {
		fmt.Printf(err.Error())
		fmt.Println("解析listener外部出口信息错误")
		return "", 0, err
	}

	err = listener.SetDeadline(time.Now().Add(time.Duration(99999 * time.Hour)))
	if err != nil {
		fmt.Printf("%s", err.Error())
		return "", 0, err
	}

	fmt.Println("我的公网IP:", strings.Split(ipPort, ":")[0])
	fmt.Println("内网的的出口端口:", port)
	return ip, port, err
}

//获取一个随机UDP Dial的内部ip，端口，外部ip端口
func GetDialIpPort(token *crypto.TokenClaims) (localAddr net.Addr, externalIp string, externalPort int, err error) {
	udpaddr, err := net.ResolveUDPAddr("udp", token.Host+":"+strconv.Itoa(token.P2PApiPort))
	//udpaddr, err := net.ResolveUDPAddr("udp", "tencent-shanghai-v1.host.nat-cloud.com:34321")
	if err != nil {
		return nil, "", 0, err
	}
	udpconn, err := net.DialUDP("udp", nil, udpaddr)
	defer udpconn.Close()
	if err != nil {
		fmt.Println(err.Error())
		return nil, "", 0, err
	}
	err = udpconn.SetDeadline(time.Now().Add(time.Duration(3 * time.Second)))
	if err != nil {
		return nil, "", 0, err
	}
	_, err = udpconn.Write([]byte("getIpPort"))
	if err != nil {
		return nil, "", 0, err
	}
	data := make([]byte, 256)
	n, err := udpconn.Read(data)
	if err != nil {
		return nil, "", 0, err
	}
	ipPort := string(data[:n])
	ip := strings.Split(ipPort, ":")[0]
	port, err := strconv.Atoi(strings.Split(ipPort, ":")[1])
	if err != nil {
		return nil, "", 0, err
	}
	//return strings.Split(udpconn.LocalAddr().String(), ":")[0]
	localAddr = udpconn.LocalAddr()
	return localAddr, ip, port, nil
}

func GetP2PListener(token *crypto.TokenClaims) (localIps string, localPort int, externalIp string, externalPort int, listener *net.UDPConn, err error) {
	localIps = GetIntranetIp()
	//localPort = randint.GenerateRangeNum(10000, 60000)
	localPort = 0
	listener, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: localPort})
	if err != nil {
		return
	}
	//获取监听的端口的外部ip和端口
	externalIp, externalPort, err = GetExternalIpPort(listener, token)
	if err != nil {
		fmt.Println(err)
		return
	}
	localPort = listener.LocalAddr().(*net.UDPAddr).Port
	return
}

//client通过指定listener发送数据到explorer指定的p2p请求地址
func SendPackToPeer(listener *net.UDPConn, ctrlmMsg *models.ReqNewP2PCtrl) {
	fmt.Println("发送包到远程：", ctrlmMsg.ExternalIp, ctrlmMsg.ExternalPort)
	//发送5次防止丢包，稳妥点
	for i := 1; i <= 5; i++ {
		listener.WriteToUDP([]byte("packFromPeer"), &net.UDPAddr{
			IP:   net.ParseIP(ctrlmMsg.ExternalIp),
			Port: ctrlmMsg.ExternalPort,
		})
		time.Sleep(time.Millisecond * 100)
	}
}
