package server

import (
	"net/http"
)

type ResponseBodyFile struct {
	Response
	file string
}

func (w *ResponseBodyFile) WriteTo(res http.ResponseWriter) {
	responseHeadersWriter.writeHeaders(w.headers, res)
	res.WriteHeader(w.statusCode)

	fileUtils.read(w.file, func(data *[]byte) {
		res.Write(*data)
	})
}

type ResponseBodyFileCachable struct {
	ResponseBodyFile
	cache    []byte
	isCached bool
}

func (w *ResponseBodyFileCachable) WriteTo(res http.ResponseWriter) {
	responseHeadersWriter.writeHeaders(w.headers, res)
	res.WriteHeader(w.statusCode)

	if w.isCached {
		res.Write(w.cache)
		return
	}

	fileUtils.read(w.file, func(data *[]byte) {
		w.cache = append(w.cache, *data...)
	})

	w.isCached = true
	res.Write(w.cache)
}
