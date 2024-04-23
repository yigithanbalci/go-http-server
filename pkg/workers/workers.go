package workers

import (
	"log"
	"net"
	"runtime"

	"github.com/yigithanbalci/go-http-server/pkg/core/method"
	"github.com/yigithanbalci/go-http-server/pkg/core/request"
	"github.com/yigithanbalci/go-http-server/pkg/handler"
)

type job struct {
	c net.Conn
}

type WorkerPool struct {
	jobs     chan *job
	handlers map[handler.HandlerConfig]handler.HandlerFunc
}

func NewPool() *WorkerPool {
	return NewPoolWithSize(0)
}

func NewPoolWithSize(size int) *WorkerPool {
	if size == 0 {
		size = runtime.NumCPU() * 2
	}
	return &WorkerPool{
		jobs: make(chan *job, size),
	}
}

func (p *WorkerPool) RegisterHandlers(handlers map[handler.HandlerConfig]handler.HandlerFunc) {
	p.handlers = handlers
}

func (p *WorkerPool) QueueConn(c net.Conn) {
	p.jobs <- &job{
		c: c,
	}
}

func (p *WorkerPool) FireWorkers() {
	for i := 1; i < runtime.NumCPU()*2; i++ {
		go worker(p.jobs, p.handlers)
	}
}

func worker(jobs <-chan *job, handlers map[handler.HandlerConfig]handler.HandlerFunc) {
	for j := range jobs {
		handleConn(j.c, handlers)
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

func handleConn(c net.Conn, handlers map[handler.HandlerConfig]handler.HandlerFunc) {
	defer c.Close()
	reqRaw, err := readFromConn(c)
	if err != nil {
		log.Printf("Error occured reading from conn: %+v", err)
	}

	req := request.ParseReq(reqRaw)
	conf := handler.HandlerConfig{
		M:    method.Method(req.Method),
		Path: req.Path,
	}
	f := handlers[conf]
	f(req)
}
