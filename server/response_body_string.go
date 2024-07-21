package server

import "net/http"

type ResponseBodyString struct {
	response []byte
	headers  map[string][]string
}

func (w *ResponseBodyString) WriteTo(res http.ResponseWriter) {
	responseHeadersWriter.writeHeaders(w.headers, res)

	res.Write(w.response)
}
