//go:build ios

package service

import (
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	"log"
	"net"
)

func GetSystemStatus(stream net.Conn, service *models.NewService) error {
	err := msg.WriteMsg(stream, &models.JsonResponse{Code: 1, Msg: "Success", Result: string("Not Support")})
	if err != nil {
		log.Println("写消息错误：")
		log.Println(err.Error())
	}
	return err
}
