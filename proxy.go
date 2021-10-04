package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"io/ioutil"
	"regexp"
	"strings"
)

type MethodType int

const (
	Method_GET  MethodType = 0
	Method_POST MethodType = 1
)

func (m MethodType) String() string {
	switch (m) {
	case Method_GET: return "GET"
	case Method_POST: return "POST"
	default: return "UNKNOWN"
	}
}

type Request struct {
	URL         string
	Method      string
	HTTPVersion string
	UserAgent   string
	Host        string
	Body        string
	Headers     map[string]string
}

func Proxy(server Server) {
	port := server.server_name + ":" + server.listen
	listen, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("listen failed, err: %v\n", err)
		return
	}
	
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("accept failed, err: %v\n", err)
			continue
		}
		go Process(conn, server)
	}
}

func Process(conn net.Conn, server Server) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
    var buf [128]byte
    n, err := reader.Read(buf[:])
    if err != nil {
        fmt.Printf("read from conn failed, err:%v\n", err)
        return
    }

    recv := string(buf[:n])
    fmt.Printf("收到的数据：\n%v\n", recv)

	request := ParseRequest(recv)
	matchLocation := MatchRequest(request, server)

	if matchLocation.category == Location_REVERSE {
		buf, n := ReverseProxy(request, matchLocation)
		if n != 0 {
			_, err = conn.Write(buf)
			if err != nil {
				fmt.Printf("write from conn failed, err:%v\n", err)
			}
		} else {
			fmt.Printf("reverse read empty")
		}
	} else {
		response := StaticProxy(request, matchLocation)
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Printf("write from conn failed, err:%v\n", err)
		}
	}
}

func ParseRequest(recv string) Request {
	request := Request{
		URL: "",
		Method: "",
		HTTPVersion: "",
		UserAgent: "",
		Host: "",
		Body: "",
		Headers: map[string]string{},
	}

	lines := strings.Split(recv, "\n")
	line := lines[0]
	params := strings.Split(line, " ")
	request.Method = params[0]
	request.URL = params[1]
	request.HTTPVersion = params[2]
	for _, line = range lines[1:] {
		params = strings.Split(line, " ")
		if len(params) != 2 {
			continue
		}
		request.Headers[params[0][:len(params)-1]] = params[1]
	}

	return request
}

func MatchRequest(request Request, server Server) Location {
	url := request.URL
	ifNormal := false
	ifRegular := false
	prefixLen := 0
	prefixLocation := Location{}
	regularLocation := Location{}
	for _, location := range server.locations {
		if location.match_pattern == Match_EAXCT && location.url == url {
			return location
		}

		if location.match_pattern == Match_NONE {
			if strings.HasPrefix(url, location.url) {
				if len(location.url) > prefixLen {
					prefixLen = len(location.url)
					prefixLocation = location
					ifNormal = false
				}
			}
		}

		if location.match_pattern == Match_NORMAL {
			if strings.HasPrefix(url, location.url) {
				if len(location.url) > prefixLen {
					prefixLen = len(location.url)
					prefixLocation = location
					ifNormal = true
				}
			}
		}

		if location.match_pattern == Match_REGULAR && !ifRegular {
			r := regexp.MustCompile(location.url)
			if len(r.FindString(url)) != 0 {
				regularLocation = location
				ifRegular = true
			}
		}

		if location.match_pattern == Match_REGULAR_NOCASE && !ifRegular {
			r := regexp.MustCompile("(?i)" + location.url)
			if len(r.FindString(url)) != 0 {
				regularLocation = location
				ifRegular = true
			}
		}
	}

	if ifNormal {
		return prefixLocation
	}
	if ifRegular {
		return regularLocation
	}
	return prefixLocation
}

func ReverseProxy(request Request, location Location) ([]byte, int) {
	reverseRequest := GenerateReverseRequest(request, location)
	conn, err := net.Dial("tcp", location.proxy_pass)
	if err != nil {
		fmt.Printf("reverse conn server failed, err:%v\n", err)
        return []byte{}, 0
	}

	_, err = conn.Write([]byte(reverseRequest))
	if err != nil {
		fmt.Printf("reverse send failed, err:%v\n", err)
		return []byte{}, 0
	}

	var buf []byte
	n, err := conn.Read(buf[:])
	if err != nil {
		fmt.Printf("reverse read failed:%v\n", err)
		return []byte{}, 0
	}

	return buf, n
}

func GenerateReverseRequest(request Request, location Location) string {
	for k, v := range location.proxy_set_header {
		request.Headers[k] = v
	}
	data := ""
	data += request.Method + " " + request.HTTPVersion + "\r\n"
	for k, v := range request.Headers {
		data += k + ":" + " " + v + "\r\n"
	}
	data += "\r\n\r\n"
	return data
}

func StaticProxy(request Request, location Location) string {
	filepath := ""
	paths := strings.Split(location.root, "/")
	for _, path := range paths {
		filepath += "/"
		filepath += path
	}
	paths = strings.Split(request.URL, "/")
	length := len(strings.Split(location.url, "/"))
	if length == len(paths) {
		filepath += "/" + location.index
	} else {
		for i, path := range paths {
			if i >= length {
				filepath += "/"
				filepath += path
			}
		}
	}

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("open file error")
	}
	defer file.Close()
	
	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("read file error")
	}

	staticResponse := GenerateStaticResponse()
	staticResponse += string(content)
	return staticResponse
}

func GenerateStaticResponse() string {
	data := ""
	data += "HTTP/1.1 200 OK\r\n"
	data += "Content-Type: image/png\r\n"
	data += "\r\n\r\n"
	return data
}