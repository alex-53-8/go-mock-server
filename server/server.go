package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Server struct {
	server *http.ServeMux
	model  *Model
}

func (s *Server) Listen() error {
	address := fmt.Sprintf(":%d", s.model.Port)
	log.Println("Listening", address)

	return http.ListenAndServe(address, s.server)
}

const responseInFilePrefix = "file:"

func createHandler(model *Model, m *Mapping, server *http.ServeMux) {
	var writer ResponseBodyWriter = nil

	if strings.HasPrefix(m.Response, responseInFilePrefix) {
		log.Println("Creating a file response writer")
		file := m.Response[len(responseInFilePrefix):]
		writer = &ResponseBodyWriterFile{
			cacheableStatus: getCacheableState(file, model.Cache.MaxItemSize),
			headers:         m.Headers,
			file:            file}
	} else {
		log.Println("Creating a string response writer")
		writer = &ResponseBodyWriterString{
			headers:  m.Headers,
			response: []byte(m.Response)}
	}

	handler := func(res http.ResponseWriter, req *http.Request) {
		log.Println(req.Method, req.RequestURI)
		writer.WriteTo(res)
	}

	if m.Method == nil || len(m.Method) == 0 {
		log.Println("Creating all methods handler for path", m.Path)
		server.HandleFunc(m.Path, handler)
	} else {
		log.Println("Creating mapping for", m.Method, "methods handler for path", m.Path)
		for _, method := range m.Method {
			server.HandleFunc(method+" "+m.Path, handler)
		}
	}
}

func NewServer(model *Model) *Server {
	var srv *http.ServeMux = http.NewServeMux()

	for _, mapping := range model.Mappings {
		createHandler(model, &mapping, srv)
	}

	return &Server{server: srv, model: model}
}
