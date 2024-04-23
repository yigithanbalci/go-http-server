package handler

import (
	"github.com/yigithanbalci/go-http-server/pkg/core/method"
	"github.com/yigithanbalci/go-http-server/pkg/core/request"
	"github.com/yigithanbalci/go-http-server/pkg/core/response"
)

type HandlerFunc func(*request.Request) *response.Response

type HandlerConfig struct {
	M    method.Method
	Path string
}
