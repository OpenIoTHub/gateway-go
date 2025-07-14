package chans

import "github.com/OpenIoTHub/gateway-go/v2/models"

// ClientTaskChan 传入访问者的ipv6监听ip+端口，任务从本chan接受配置创建客户端
var ClientTaskChan = make(chan models.Ipv6ClientHandleTask)
