package workers

import (
	"fmt"
	"log"
	"net"
	"runtime"

	"github.com/yigithanbalci/go-http-server/pkg/request"
)

func init() {
	go fireWorkers()
}

type job struct {
	c net.Conn
}

var jobs = make(chan *job, runtime.NumCPU()*2)

func QueueConn(c net.Conn) {
	jobs <- &job{
		c: c,
	}
}

func fireWorkers() {
	for i := 1; i < runtime.NumCPU()*2; i++ {
		go worker(jobs)
	}
}

func worker(jobs <-chan *job) {
	for j := range jobs {
		handleConn(j.c)
	}
}

var requestRaw = make([]byte, 1024)

func readFromConn(r net.Conn) ([]byte, error) {
	n, err := r.Read(requestRaw)
	if err != nil {
		return nil, err
	}

	return requestRaw[:n], nil
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

func handleConn(c net.Conn) {
	defer c.Close()
	reqRaw, err := readFromConn(c)
	if err != nil {
		log.Printf("Error occured reading from conn: %+v", err)
	}

	req := request.ParseReq(reqRaw)
	// var bytesMsg []byte
	// re := response.Response{
	// 	HttpVer:       req.HttpVer,
	// 	Status:        "404 NOT FOUND",
	// 	ContentType:   "",
	// 	ContentLength: "",
	// 	Body:          "",
	// }

	// if strings.Contains(req.path, "/echo") {
	// 	echoHandler(req, &re)
	// }
	// if req.path == "/" {
	// 	baseHandler(&re)
	// }
	// if req.path == "/user-agent" {
	// 	userAgentHandler(req, &re)
	// }
	// if strings.Contains(req.path, "/files") {
	// 	fileHandler(req, &re)
	// }
	// bytesMsg = parseResp(&re)
	// fmt.Println("msg:", string(bytesMsg))
	// err = writeToConn(c, bytesMsg)
	// fmt.Println("written")
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
