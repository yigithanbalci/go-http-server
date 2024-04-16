package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const separator = "\r\n"

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

type resp struct {
	httpVer       string
	status        string
	contentType   string
	contentLength string
	body          string
}

func parseResp(re *resp) []byte {
	respStr := re.httpVer + " "
	respStr = respStr + re.status + separator
	if re.contentType != "" {
		respStr = respStr + "Content-Type: " + re.contentType + separator
	}
	if re.contentLength != "" {
		respStr = respStr + "Content-Length: " + re.contentLength + separator
	}
	if re.body != "" {
		respStr = respStr + separator
		respStr = respStr + re.body
	}
	respStr = respStr + separator + separator
	fmt.Printf("respstr:%s", respStr)
	fmt.Printf("resp: %+v", re)
	return []byte(respStr)
}

type req struct {
	method        string
	path          string
	httpVer       string
	userAgent     string
	contentLength string
	body          string
}

func parseReq(msg []byte) *req {
	smsg := string(msg)
	contents := strings.Split(smsg, separator)
	request := &req{
		method:        "",
		path:          "",
		httpVer:       "",
		contentLength: "",
		body:          "",
	}
	for i, e := range contents {
		if i == 0 {
			startLine := strings.Split(contents[0], " ")
			request.method = startLine[0]
			request.path = startLine[1]
			request.httpVer = startLine[2]
		} else {
			if strings.Contains(e, "User-Agent:") {
				agent := strings.SplitAfter(e, "User-Agent: ")
				request.userAgent = agent[1]
			} else if strings.Contains(e, "Content-Length:") {
				contentLength := strings.SplitAfter(e, "Content-Length: ")
				request.contentLength = contentLength[1]
			}
		}
	}

	if request.method == "POST" && request.contentLength != "" {
		request.body = contents[len(contents)-1]
	}
	return request
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
