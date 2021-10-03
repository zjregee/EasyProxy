package main

type LocationType int

type MatchType int

const (
	Location_REVERSE LocationType = 0
	Location_STATIC  LocationType = 1
)

const (
	Match_NONE           MatchType = 0
	Match_EAXCT          MatchType = 1
	Match_REGULAR        MatchType = 2
	Match_REGULAR_NOCASE MatchType = 3
	Match_NORMAL         MatchType = 4
)

const (
	Access_Log_DEFAULT string = "access.log"
	Error_Log_DEFAULT  string = "error.log"
)

type Server struct {
	listen      string
	server_name string
	locations   []Location
	error_log   string
	access_log  string
}

type Location struct {
	category         LocationType             
	proxy_set_header map[string]string
	proxy_pass       string
	root             string
	index            string
	match_pattern    MatchType
	url              string
}

type Config struct {
	servers []Server
}

var ProxyConfig  Config

func InitConfig() {
	server1 := Server {
		listen: "80",
		server_name: "localhost1",
		error_log: "/mnt/log/nginx/localhost1/error.log",
		access_log: "/mnt/log/nginx/localhost1/access.log",
	}
	server2 := Server {
		listen: "80",
		server_name: "localhost2",
		error_log: "/mnt/log/nginx/localhost2/error.log",
		access_log: "/mnt/log/nginx/localhost2/access.log",
	}

	location1 := Location {
		category: Location_REVERSE,
		proxy_pass: "127.0.0.1:6550",
		match_pattern: Match_NORMAL,
		url: "/api/v1",
	}

	location1.proxy_set_header = map[string]string{}
	location1.proxy_set_header["Host"] = "host"

	location2 := Location {
		category: Location_STATIC,
		root: "/mnt/var/www/localhost1",
		index: "index.html",
		match_pattern: Match_NORMAL,
		url: "/static/",
	}

	location3 := Location {
		category: Location_REVERSE,
		proxy_pass: "127.0.0.1:6551",
		match_pattern: Match_NORMAL,
		url: "/api/v1",
	}

	location3.proxy_set_header = map[string]string{}
	location3.proxy_set_header["Host"] = "host"

	location4 := Location {
		category: Location_STATIC,
		root: "/mnt/var/www/localhost2",
		index: "index.html",
		match_pattern: Match_NORMAL,
		url: "/static/",
	}

	server1.locations = []Location{}
	server1.locations = append(server1.locations, location1)
	server1.locations = append(server1.locations, location2)

	server2.locations = []Location{}
	server2.locations = append(server2.locations, location3)
	server2.locations = append(server2.locations, location4)

	ProxyConfig = Config{}
	ProxyConfig.servers = []Server{}
	ProxyConfig.servers = append(ProxyConfig.servers, server1)
	ProxyConfig.servers = append(ProxyConfig.servers, server2)
}