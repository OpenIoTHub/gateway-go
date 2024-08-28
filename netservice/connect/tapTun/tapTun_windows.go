package tapTun

import (
	"errors"
	"github.com/OpenIoTHub/utils/io"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/water"
	"log"
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
	log.Println("虚拟网卡名称：", ifaceName)
	cmd := exec.Command("netsh", "interface", "ip", "set", "address",
		"name", "=", ifaceName, "source", "=", "static", "addr", "=", "192.168.69.1", "mask", "=", "255.255.255.0", "gateway", "=", "none")
	err = cmd.Run()
	if err != nil {
		return errors.New("请以管理员身份运行本软件！" + err.Error())
	}
	go io.Join(stream, ifce)
	return nil
}
