package handler

import (
	"github.com/yigithanbalci/go-http-server/pkg/core/method"
	"github.com/yigithanbalci/go-http-server/pkg/request"
	"github.com/yigithanbalci/go-http-server/pkg/response"
)

type HandlerFunc func(*request.Request) *response.Response

type HandlerConfig struct {
	M           method.Method
	ContentType string
	Path        string
}

var handlers = make(map[HandlerConfig]HandlerFunc)

func AddHandler(conf HandlerConfig, f HandlerFunc) {
	handlers[conf] = f
}
