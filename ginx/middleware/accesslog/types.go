package accesslog

import (
	"github.com/gin-gonic/gin"
)

// AccessLog 你可以打印很多的信息，根据需要自己加
type AccessLog struct {
	PID      string `json:"pid"`
	Referer  string `json:"referer"`
	Protocol string `json:"protocol"`
	Port     string `json:"port"`
	IP       string `json:"ip"`
	IPs      string `json:"ips"`
	Host     string `json:"host"`
	ClientIP string `json:"client_ip"`
	URL      string `json:"url"`
	UA       string `json:"ua"`

	Method     string `json:"method"`
	Path       string `json:"path"`
	ReqBody    string `json:"req_body"`
	Duration   string `json:"duration"`
	StatusCode int    `json:"status_code"`
	RespBody   string `json:"resp_body"`
}

type responseWriter struct {
	al *AccessLog
	gin.ResponseWriter
	maxLength int64
}

func (r responseWriter) WriteHeader(statusCode int) {
	r.al.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r responseWriter) Write(data []byte) (int, error) {
	curLen := int64(len(data))
	if curLen >= r.maxLength {
		data = data[:r.maxLength]
	}
	r.al.RespBody = string(data)
	return r.ResponseWriter.Write(data)
}

func (r responseWriter) WriteString(data string) (int, error) {
	curLen := int64(len(data))
	if curLen >= r.maxLength {
		data = data[:r.maxLength]
	}
	r.al.RespBody = data
	return r.ResponseWriter.WriteString(data)
}
