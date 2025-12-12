package docker

import (
	"fmt"
	"testing"
)

func TestGetContainersInfo(t *testing.T) {
	got, err := GetContainersInfo()
	if err != nil {
		t.Fatal(err)
	}
	for _, info := range got.Items {
		t.Log(info)
		fmt.Printf("%+v", info)
		//t.Log(info.Ports)
		//t.Log(info.Status)
		//t.Log(info.State)
	}
}

func TestGetContainersServices(t *testing.T) {
	clr := GetContainersServices()
	for _, info := range clr {
		t.Log(info)
	}
}
