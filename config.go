package main

type Server struct {
	listen      string
	server_name string
	locations   []Location
	error_log   string
	access_log  string
}

type Location struct {
	proxy_set_header map[string]string
	proxy_pass       string
	root             string
	index            string
}

type Config struct {
	servers []Server
}

func (c *Config) New() {
	// 读取配置文件
	// 初始化Config
}