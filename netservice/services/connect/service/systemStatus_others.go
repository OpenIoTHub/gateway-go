//go:build !ios

package service

import (
	"encoding/json"
	"github.com/OpenIoTHub/utils/v2/models"
	"github.com/OpenIoTHub/utils/v2/msg"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	sysnet "github.com/shirou/gopsutil/v3/net"

	"log"
	"net"
)

func GetSystemStatus(stream net.Conn, service *models.NewService) error {
	statMap := make(map[string]interface{})
	// 获取主机相关信息
	hostInfo, _ := host.Info()
	hostMap := make(map[string]interface{})
	hostMap["uptime"] = hostInfo.Uptime                   //运行时间
	hostMap["bootTime"] = hostInfo.BootTime               //启动时间
	hostMap["procs"] = hostInfo.Procs                     //进程数
	hostMap["os"] = hostInfo.OS                           //操作系统
	hostMap["platform"] = hostInfo.Platform               //平台
	hostMap["platformVersion"] = hostInfo.PlatformVersion //平台版本
	hostMap["kernelArch"] = hostInfo.KernelArch           //内核
	hostMap["kernelVersion"] = hostInfo.KernelVersion     //内核版本
	statMap["hosts"] = hostMap

	// 获取内存信息
	memInfo, _ := mem.VirtualMemory()
	memMap := make(map[string]interface{})
	memMap["total"] = memInfo.Total             //总内存
	memMap["available"] = memInfo.Available     //可用内存
	memMap["used"] = memInfo.Used               //已使用内存
	memMap["free"] = memInfo.Free               //剩余内存
	memMap["usedPercent"] = memInfo.UsedPercent //百分比
	memMap["buffers"] = memInfo.Buffers         //缓存
	memMap["shared"] = memInfo.Shared           //共享内存
	memMap["cached"] = memInfo.Cached           //缓冲区
	statMap["mems"] = memMap

	// 获取CPU信息
	cpuInfo, _ := cpu.Info()
	var cpuMapArr []map[string]interface{}
	for _, c := range cpuInfo {
		cpuMap := make(map[string]interface{})
		cpuMap["cpu"] = c.CPU + 1         //第几个CPU 从0开始的
		cpuMap["cores"] = c.Cores         //CPU的核数
		cpuMap["modelName"] = c.ModelName //CPU类型
		cpuMapArr = append(cpuMapArr, cpuMap)
	}
	statMap["cpus"] = cpuMapArr

	// 获取IO信息
	ioInfo, _ := sysnet.IOCounters(false)
	var ioMapArr []map[string]interface{}
	for _, i := range ioInfo {
		ioMap := make(map[string]interface{})
		ioMap["ioName"] = i.Name             //网口名
		ioMap["bytesSent"] = i.BytesSent     //发送字节数
		ioMap["bytesRecv"] = i.BytesRecv     //接收字节数
		ioMap["packetsSent"] = i.PacketsSent //发送的数据包数
		ioMap["packetsRecv"] = i.PacketsRecv //接收的数据包数
		ioMapArr = append(ioMapArr, ioMap)
	}
	statMap["ios"] = ioMapArr

	// 获取磁盘信息
	partitions, _ := disk.Partitions(false)
	var diskMapArr []map[string]interface{}
	for _, partition := range partitions {
		diskMap := make(map[string]interface{})
		usage, _ := disk.Usage(partition.Mountpoint)
		diskMap["disk"] = partition.Mountpoint     //第几块磁盘
		diskMap["total"] = usage.Total             //总大小
		diskMap["free"] = usage.Free               //剩余空间
		diskMap["used"] = usage.Used               //已使用空间
		diskMap["usedPercent"] = usage.UsedPercent //百分比
		diskMapArr = append(diskMapArr, diskMap)
	}
	statMap["disks"] = diskMapArr
	rstByte, err := json.Marshal(statMap)
	if err != nil {
		log.Println("json.Marshal(statMap)：")
		log.Println(err.Error())
	}
	//log.Println(string(rstByte))
	err = msg.WriteMsg(stream, &models.JsonResponse{Code: 0, Msg: "Success", Result: string(rstByte)})
	if err != nil {
		log.Println("写消息错误：")
		log.Println(err.Error())
	}
	return err
}
