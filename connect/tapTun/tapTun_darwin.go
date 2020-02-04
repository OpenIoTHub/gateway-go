package tapTun

import (
	"fmt"
	"github.com/OpenIoTHub/utils/io"
	"github.com/OpenIoTHub/utils/models"
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
	fmt.Println("ifaceName", ifaceName)
	cmd := exec.Command("ifconfig", ifaceName, "192.168.69.1", "netmask", "255.255.255.0", "broadcast", "192.168.69.255", "up")
	err = cmd.Run()
	if err != nil {
		return err
	}
	go io.Join(stream, ifce)
	return nil
}
