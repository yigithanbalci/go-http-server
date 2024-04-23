package workers

import (
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
	_, err := w.Write(msg)
	if err != nil {
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
}
