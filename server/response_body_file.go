package server

import (
	"net/http"
)

type ResponseBodyFile struct {
	file    string
	headers map[string][]string
}

func (w *ResponseBodyFile) WriteTo(res http.ResponseWriter) {
	responseHeadersWriter.writeHeaders(w.headers, res)

	fileUtils.read(w.file, func(data *[]byte) {
		res.Write(*data)
	})
}

type ResponseBodyFileCachable struct {
	file     string
	cache    []byte
	isCached bool
	headers  map[string][]string
}

func (w *ResponseBodyFileCachable) WriteTo(res http.ResponseWriter) {
	responseHeadersWriter.writeHeaders(w.headers, res)

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
