package main

import "C"

import (
	client "github.com/OpenIoTHub/gateway-go/client"
)

var (
	version = "dev"
	commit  = ""
	date    = "time"
	builtBy = ""
)

func main() {

}

//export Run
func Run() {
	go client.Run()
}
