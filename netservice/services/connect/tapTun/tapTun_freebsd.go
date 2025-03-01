package tapTun

import (
	"github.com/OpenIoTHub/utils/models"
	"net"
)

func NewTun(stream net.Conn, service *models.NewService) error {

	return NewTap(stream, service)
}

func NewTap(stream net.Conn, service *models.NewService) error {
	return nil
}
