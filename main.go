package main

func main() {
	InitConfig()
	logs := []string{}
	for _, server := range ProxyConfig.servers {
		if server.access_log != "" {
			logs = append(logs, server.access_log)
		} else {
			logs = append(logs, Access_Log_DEFAULT)
		}
		if server.error_log != "" {
			logs = append(logs, server.error_log)
		} else {
			logs = append(logs, Error_Log_DEFAULT)
		}
	}

	InitLog(logs)
	defer CloseLog()
	
	for _, server := range ProxyConfig.servers {
		go Proxy(server)
	}
}