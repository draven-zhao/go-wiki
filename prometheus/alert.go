package prometheus

func BuildQuery(a Alert) (string, string) {
	alertName := Normalize(a.Labels["alertname"])
	return alertName, a.Labels["instance"]
}
func Normalize(name string) string {
	switch name {
	case "node_disk_utilization_rate":
		return "磁盘使用率过高"
	case "node_memory_utilization_rate":
		return "内存使用率过高"
	case "node_arp_entries":
		return "node arp entries"
	case "dcos_service_number_master":
		return "dcos master服务数量"
	case "dcos_service_state_slave":
		return "dcos slave服务状态"
	case "node_load_average_1m":
		return "节点负载1m平均值"
	case "node_mesos_slave_single_node_status":
		return "mesos slave节点状态"
	case "node_status":
		return "节点状态"
	case "dcos_service_state_master":
		return "dcos master 服务状态"
	case "node_docker_status":
		return "docker 状态"
	case "node_exporter_status":
		return "node exporter 状态"
	case "node_ping_loss":
		return "节点ping 丢包"
	case "node_disk_io_time":
		return "磁盘io iops"
	case "monitor_single_exporter_status":
		return "监控节点 exporter 状态"
	case "node_ping_timeout":
		return "节点ping超时"
	}
	return name
	// todo 根据alertname做返回和召回，不要全部返回
}
