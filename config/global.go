package config

// Ipv6ListenTcpHandlePort ipv6监听访问端直接的连接，
// 处理测试链接，测试链接显示可以连通后面有请求直接使用，当连接不通则使用其他连接
// 这个端口在启动时开启，开启后将实际端口号存到这里，当接到访问者请求之后将该端口和ipv6地址发送到访问者
var Ipv6ListenTcpHandlePort int = 0
