package tapTun

import (
	"fmt"
	"github.com/OpenIoTHub/utils/v2/io"
	"github.com/OpenIoTHub/utils/v2/models"
	"github.com/OpenIoTHub/water"
	"net"
	"os/exec"
	"time"
)

func NewTun(stream net.Conn, service *models.NewService) error {

	return NewTap(stream, service)
}

func NewTap(stream net.Conn, service *models.NewService) error {
	config := water.Config{
		DeviceType: water.TAP,
	}
	config.Name = "OpenIoTHub" + fmt.Sprintf("-%d", time.Now().UTC().Unix())
	ifce, err := water.New(config)
	if err != nil {
		return err
	}
	cmd1 := exec.Command("ip", "addr", "add", "192.168.69.1/24", "dev", ifce.Name())
	err = cmd1.Run()
	if err != nil {
		return err
	}
	cmd2 := exec.Command("ip", "link", "set", "dev", ifce.Name(), "up")
	err = cmd2.Run()
	if err != nil {
		return err
	}
	go io.Join(stream, ifce)
	return nil
}
