//go:build ios

package service

import (
	"encoding/json"
	"log"
	"net"

	"github.com/OpenIoTHub/utils/v2/models"
	"github.com/OpenIoTHub/utils/v2/msg"
)

func GetSystemStatus(stream net.Conn, service *models.NewService) error {
	statMap := make(map[string]interface{})
	statMap["code"] = 1
	statMap["message"] = "failed"
	rstByte, err := json.Marshal(statMap)
	if err != nil {
		log.Println("json.Marshal(statMap)：")
		log.Println(err.Error())
	}
	err = msg.WriteMsg(stream, &models.JsonResponse{Code: 1, Msg: "Success", Result: string(rstByte)})
	if err != nil {
		log.Println("写消息错误：")
		log.Println(err.Error())
	}
	return err
}
