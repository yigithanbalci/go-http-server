package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func handleConn(c net.Conn) {
	defer c.Close()
	req, err := readFromConn(c)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v", req)
	request := parseReq(req)
	fmt.Printf("request:%+v", request)
	var bytesMsg []byte
	re := resp{
		httpVer:       request.httpVer,
		status:        "404 NOT FOUND",
		contentType:   "",
		contentLength: "",
		body:          "",
	}

	if strings.Contains(request.path, "/echo") {
		echoHandler(request, &re)
	}
	if request.path == "/" {
		baseHandler(&re)
	}
	if request.path == "/user-agent" {
		userAgentHandler(request, &re)
	}
	if strings.Contains(request.path, "/files") {
		fileHandler(request, &re)
	}
	bytesMsg = parseResp(&re)
	fmt.Println("msg:", string(bytesMsg))
	err = writeToConn(c, bytesMsg)
	fmt.Println("written")
	if err != nil {
		fmt.Println(err)
	}
}

func baseHandler(re *resp) {
	re.status = "200 OK"
}

func userAgentHandler(request *req, re *resp) {
	re.status = "200 OK"
	re.contentType = "text/plain"
	re.contentLength = strconv.Itoa(len(request.userAgent))
	re.body = request.userAgent
}

func fileHandler(request *req, re *resp) {
	file := strings.SplitAfter(request.path, "/files/")
	switch request.method {
	case "GET":
		cont, err := os.ReadFile(fileDir + file[1])
		if err != nil {
			re.status = "404 NOT FOUND"
			return
		}

		re.status = "200 OK"
		re.contentType = "application/octet-stream"
		re.contentLength = strconv.Itoa(len(cont))
		re.body = string(cont)
	case "POST":
		err := os.WriteFile(fileDir+file[1], []byte(request.body), 0644)
		if err != nil {
			re.status = "500 INTERNAL SERVER ERROR"
			return
		}
		re.status = "201 CREATED"
	}
}

func echoHandler(request *req, re *resp) {
	strs := strings.Split(request.path, "/echo/")

	var subPath string
	if len(strs) == 2 {
		subPath = strs[1]
		if subPath != "" {
			re.status = "200 OK"
			re.contentType = "text/plain"
			re.contentLength = strconv.Itoa(len(subPath))
			re.body = subPath
		}
	}
}

func writeToConn(w net.Conn, msg []byte) error {
	n, err := w.Write(msg)
	if err != nil {
		fmt.Println(n)
		fmt.Println("err: ", err)
		return err
	}

	return nil
}

var request = make([]byte, 1024)

func readFromConn(r net.Conn) ([]byte, error) {
	n, err := r.Read(request)
	if err != nil {
		fmt.Println(n)
		return nil, err
	}

	return request[:n], nil
}

var fileDir string

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	flag.StringVar(&fileDir, "directory", "/", "file directory")
	flag.Parse()
	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConn(conn)
	}
}
