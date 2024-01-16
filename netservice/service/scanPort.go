package service

import (
	"encoding/json"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	"log"
	"net"
)

func loop(startport, endport int, inport chan int) {
	for i := startport; i <= endport; i++ {
		inport <- i
	}
}

func scanner(inport, outport, out chan int, ip net.IP, endport int) {
	for {
		in := <-inport
		tcpaddr := &net.TCPAddr{IP: ip, Port: in}
		conn, err := net.DialTCP("tcp", nil, tcpaddr)
		if err != nil {
			outport <- 0
		} else {
			outport <- in
			conn.Close()
		}
		if in == endport {
			out <- in
		}
	}
}

func ScanPort(stream net.Conn, service *models.NewService) error {
	//decode json
	var config *models.ScanPort
	//var rst *models.ScanPortResult
	err := json.Unmarshal([]byte(service.Config), &config)
	if err != nil {
		return err
	}
	inport := make(chan int)
	outport := make(chan int)
	out := make(chan int)
	collect := []int{}
	go loop(config.StartPort, config.EndPort, inport)
	for {
		needBreak := false
		select {
		case <-out:
			//log.Println(collect)
			needBreak = true
		default:
			ip := net.ParseIP(config.Host)
			go scanner(inport, outport, out, ip, config.EndPort)
			port := <-outport
			if port != 0 {
				collect = append(collect, port)
			}
		}
		if needBreak {
			break
		}
	}
	rstByte, err := json.Marshal(&collect)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	//log.Println(string(rstByte))
	err = msg.WriteMsg(stream, &models.JsonResponse{Code: 0, Msg: "Success", Result: string(rstByte)})
	if err != nil {
		log.Println("写消息错误：")
		log.Println(err.Error())
	}
	return err
}
