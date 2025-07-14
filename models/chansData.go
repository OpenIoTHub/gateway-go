package models

type Ipv6ClientHandleTask struct {
	RunId        string `json:"RunId"`
	Ipv6AddrIp   string `json:"Ipv6AddrIp"`
	Ipv6AddrPort int    `json:"Ipv6AddrPort"`
}
