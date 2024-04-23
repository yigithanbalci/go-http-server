package main

import (
	"fmt"
	"net"
	"os"

	"github.com/yigithanbalci/go-http-server/pkg/core/content_type"
	"github.com/yigithanbalci/go-http-server/pkg/core/method"
	"github.com/yigithanbalci/go-http-server/pkg/core/request"
	"github.com/yigithanbalci/go-http-server/pkg/core/response"
	"github.com/yigithanbalci/go-http-server/pkg/core/status"
	"github.com/yigithanbalci/go-http-server/pkg/core/version"
	"github.com/yigithanbalci/go-http-server/pkg/handler"
	"github.com/yigithanbalci/go-http-server/pkg/workers"
)

type Server struct {
	HttpVer    string
	Port       string
	handlers   map[handler.HandlerConfig]handler.HandlerFunc
	workerPool *workers.WorkerPool
}

func NewServer() *Server {
	return &Server{
		HttpVer:    version.V1_1,
		Port:       "8080",
		handlers:   make(map[handler.HandlerConfig]handler.HandlerFunc),
		workerPool: workers.NewPool(),
	}
}

func (s *Server) AddHandler(m method.Method, path string, f handler.HandlerFunc) {
	conf := handler.HandlerConfig{
		M:    m,
		Path: path,
	}
	s.handlers[conf] = f
}

func (s *Server) Run() {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", "0.0.0.0", s.Port))
	if err != nil {
		fmt.Println("Failed to bind to port ", s.Port)
		os.Exit(1)
	}
	defer l.Close()
	s.workerPool.RegisterHandlers(s.handlers)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		s.workerPool.QueueConn(conn)
	}
}

func main() {
	s := NewServer()
	s.AddHandler(method.GET, "/", func(r *request.Request) *response.Response {
		return &response.Response{
			HttpVer:       r.HttpVer,
			Status:        status.OK,
			ContentType:   content_type.APPLICATION_JSON,
			ContentLength: string(len([]byte("Helloweeb"))),
			Body:          "Helloweeb",
		}
	})
	s.Run()
}
