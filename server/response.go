package server

import (
	"net/http"
)

type ResponseBody interface {
	WriteTo(res http.ResponseWriter)
}

type ResponseHeaders interface {
	writeHeaders(headers map[string][]string, res http.ResponseWriter)
}

type ResponseHeadersWriter struct {
}

var responseHeadersWriter ResponseHeaders = &ResponseHeadersWriter{}

func (r *ResponseHeadersWriter) writeHeaders(headers map[string][]string, res http.ResponseWriter) {
	for name, values := range headers {
		for _, value := range values {
			res.Header().Add(name, value)
		}
	}
}
