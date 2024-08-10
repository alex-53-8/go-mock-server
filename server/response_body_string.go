package server

import "net/http"

type ResponseBodyString struct {
	Response
	response []byte
}

func (w *ResponseBodyString) WriteTo(res http.ResponseWriter) {
	responseHeadersWriter.writeHeaders(w.headers, res)
	res.WriteHeader(w.statusCode)
	res.Write(w.response)
}
