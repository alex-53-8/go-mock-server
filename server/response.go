package server

import (
	"io"
	"log"
	"net/http"
	"os"
)

type ResponseBodyWriter interface {
	WriteTo(res http.ResponseWriter)
}

type cachable int

const canBeCached cachable = 1
const cannotBeCached cachable = 2

type ResponseBodyWriterFile struct {
	file            string
	cache           []byte
	isCached        bool
	cacheableStatus cachable
	headers         []Header
}

const bufferSize = 1024

func readAndWriteToResponse(filename string, res http.ResponseWriter) {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModeDevice)
	if err != nil {
		log.Println("cannot read response file: ", filename, err.Error())

	}

	defer file.Close()
	buf := make([]byte, bufferSize)

	for {
		n, err := file.Read(buf)
		if err == io.EOF || n <= 0 {
			break
		}
		if err == nil {
			res.Write(buf[:n])
		}
	}
}

func readAllBytes(filename string) []byte {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModeDevice)
	if err != nil {
		log.Println("cannot read response file: ", filename, err.Error())

	}

	defer file.Close()
	buf := make([]byte, bufferSize)
	response := make([]byte, 0)

	for {
		n, err := file.Read(buf)
		if err == io.EOF || n <= 0 {
			break
		}
		if err == nil {
			response = append(response, buf[:n]...)
		}
	}

	return response
}

func getCacheableState(filename string, maxCacheableSize int64) cachable {
	info, err := os.Stat(filename)

	if err != nil {
		log.Println("cannot get file stat: ", filename, err.Error())
		return cannotBeCached
	}

	size := info.Size()
	if size < maxCacheableSize {
		return canBeCached
	} else {
		return cannotBeCached
	}
}

func (w *ResponseBodyWriterFile) WriteTo(res http.ResponseWriter) {
	for _, header := range w.headers {
		res.Header().Add(header.Key, header.Value)
	}

	if w.isCached {
		res.Write(w.cache)
		return
	}

	if w.cacheableStatus == canBeCached {
		w.cache = readAllBytes(w.file)
		w.isCached = true
		res.Write(w.cache)
	} else {
		readAndWriteToResponse(w.file, res)
	}
}

type ResponseBodyWriterString struct {
	response []byte
	headers  []Header
}

func (w *ResponseBodyWriterString) WriteTo(res http.ResponseWriter) {
	for _, header := range w.headers {
		res.Header().Add(header.Key, header.Value)
	}
	res.Write(w.response)
}
