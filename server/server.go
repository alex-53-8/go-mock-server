package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Server struct {
	server *http.Server
	model  *Model
}

func (s *Server) Listen() error {
	log.Println("Listening", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shuting down...")
	return s.server.Shutdown(ctx)
}

const responseInFilePrefix = "file:"

func createWriterFileResponse(cfg *Cfg, m *Endpoint) ResponseBody {
	log.Println("Creating a file response writer")
	file := m.Response[len(responseInFilePrefix):]

	if cfg.CachingEnabled && getCacheableState(file, cfg.CacheItemMaxSize) == canBeCached {
		log.Println(file, "can be cached")
		return &ResponseBodyFileCachable{
			headers: m.Headers,
			file:    file}
	} else {
		log.Println(file, "cannot be cached, will read each time")
		return &ResponseBodyFile{
			headers: m.Headers,
			file:    file}
	}
}

func createHandler(cfg *Cfg, m *Endpoint, server *http.ServeMux) {
	var writer ResponseBody = nil

	if strings.HasPrefix(m.Response, responseInFilePrefix) {
		writer = createWriterFileResponse(cfg, m)
	} else {
		log.Println("Creating a string response writer")
		writer = &ResponseBodyString{
			headers:  m.Headers,
			response: []byte(m.Response)}
	}

	handler := func(res http.ResponseWriter, req *http.Request) {
		log.Println(req.Method, req.RequestURI)
		writer.WriteTo(res)
	}

	if m.Method == nil || len(m.Method) == 0 {
		log.Println("Mapping all methods handler for path", m.Path)
		server.HandleFunc(m.Path, handler)
	} else {
		log.Println("Mapping ", m.Method, "methods handler for path", m.Path)
		for _, method := range m.Method {
			server.HandleFunc(strings.ToUpper(method)+" "+m.Path, handler)
		}
	}
}

func NewServer(model *Model, cfg Cfg) *Server {
	var srv *http.ServeMux = http.NewServeMux()

	for _, mapping := range model.Endpoints {
		createHandler(&cfg, &mapping, srv)
	}

	address := fmt.Sprintf(":%d", model.Port)

	s := &http.Server{Addr: address, Handler: srv}
	return &Server{server: s, model: model}
}
