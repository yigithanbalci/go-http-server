package request

import (
	"strings"

	consts "github.com/yigithanbalci/go-http-server/pkg/core"
)

type Request struct {
	Method        string
	Path          string
	HttpVer       string
	UserAgent     string
	ContentLength string
	Body          string
}

func parseReq(msg []byte) *Request {
	smsg := string(msg)
	contents := strings.Split(smsg, consts.Separator)
	request := &Request{
		Method:        "",
		Path:          "",
		HttpVer:       "",
		ContentLength: "",
		Body:          "",
	}
	for i, e := range contents {
		if i == 0 {
			startLine := strings.Split(contents[0], " ")
			request.Method = startLine[0]
			request.Path = startLine[1]
			request.HttpVer = startLine[2]
		} else {
			if strings.Contains(e, "User-Agent:") {
				agent := strings.SplitAfter(e, "User-Agent: ")
				request.UserAgent = agent[1]
			} else if strings.Contains(e, "Content-Length:") {
				contentLength := strings.SplitAfter(e, "Content-Length: ")
				request.ContentLength = contentLength[1]
			}
		}
	}

	if request.Method == "POST" && request.ContentLength != "" {
		request.Body = contents[len(contents)-1]
	}
	return request
}
