package response

import (
	"fmt"

	"github.com/yigithanbalci/go-http-server/pkg/core"
)

type Response struct {
	HttpVer       string
	Status        string
	ContentType   string
	ContentLength string
	Body          string
}

func ParseResp(re *Response) []byte {
	respStr := re.HttpVer + " "
	respStr = respStr + re.Status + core.Separator
	if re.ContentType != "" {
		respStr = respStr + "Content-Type: " + re.ContentType + core.Separator
	}
	if re.ContentLength != "" {
		respStr = respStr + "Content-Length: " + re.ContentLength + core.Separator
	}
	if re.Body != "" {
		respStr = respStr + core.Separator
		respStr = respStr + re.Body
	}
	respStr = respStr + core.Separator + core.Separator
	fmt.Printf("respstr:%s", respStr)
	fmt.Printf("resp: %+v", re)
	return []byte(respStr)
}
