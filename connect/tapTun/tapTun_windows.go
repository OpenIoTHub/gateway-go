package tapTun

import (
	"errors"
	"fmt"
	"git.iotserv.com/iotserv/utils/io"
	"git.iotserv.com/iotserv/utils/models"
	"github.com/songgao/water"
	"net"
	"os/exec"
)

func NewTun(stream net.Conn, service *models.NewService) error {

	return NewTap(stream, service)
}

func NewTap(stream net.Conn, service *models.NewService) error {
	config := water.Config{
		DeviceType: water.TAP,
	}
	ifce, err := water.New(config)
	if err != nil {
		return err
	}
	ifaceName := ifce.Name()
	fmt.Println("虚拟网卡名称：", ifaceName)
	cmd := exec.Command("netsh", "interface", "ip", "set", "address",
		"name", "=", ifaceName, "source", "=", "static", "addr", "=", "192.168.69.1", "mask", "=", "255.255.255.0", "gateway", "=", "none")
	err = cmd.Run()
	if err != nil {
		return errors.New("请以管理员身份运行本软件！" + err.Error())
	}
	go io.Join(stream, ifce)
	return nil
}
