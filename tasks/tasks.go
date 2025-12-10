package tasks

func RunTasks() {
	go ipv6ServerTask()
	go ipv6ClientTask()
}
