package main

import (
	"bufio"
	"fmt"
	"net"
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

func parseRequest(recv string) (Request, error) {
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
		request.Headers[params[0][:len(params)-2]] = params[1]
	}

	return request, nil
}

func generateRequest() (string, error) {
	data := "GET / HTTP/1.1\r\n"
    data += "HOST: www.baidu.com\r\n"
    data += "connection: close\r\n"
    data += "\r\n\r\n"
	return data, nil
}

func process(conn net.Conn) {
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
	// _, _ = parseRequest(recv)

    _, err = conn.Write([]byte("ok"))
    if err != nil {
        fmt.Printf("write from conn failed, err:%v\n", err)
    }
}

func fun1() {
	fmt.Println("fun1")
	conn, err := net.Dial("tcp", "www.baidu.com:80")
    if err != nil {
        fmt.Printf("conn server failed, err:%v\n", err)
        return
    }

	s, _ := generateRequest()

	_, err = conn.Write([]byte(s))
	if err != nil {
		fmt.Printf("send failed, err:%v\n", err)
		return
	}
	
	var buf [8192*10]byte
	n, err := conn.Read(buf[:])
	if err != nil {
		fmt.Printf("read failed:%v\n", err)
		return
	}
	fmt.Printf("收到服务端回复:%v\n", string(buf[:n]))
	fmt.Println("fun2")
}

func main() {
	// listen, err := net.Listen("tcp", "127.0.0.1:1313")
	// if err != nil {
	// 	fmt.Printf("listen failed, err: %v\n", err)
	// 	return
	// }
	
	// for {
	// 	conn, err := listen.Accept()
	// 	if err != nil {
	// 		fmt.Printf("accept failed, err: %v\n", err)
	// 		continue
	// 	}
	// 	go process(conn)
	// }

	// fun1()
}