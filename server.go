package httpserver

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	// Router holds a performant router for requests
	Router *httprouter.Router

	// Middlewares holds global middlewares for the server
	Middlewares []func(next http.Handler) http.Handler
}

// New creates a server
func New(middlewares ...func(http.Handler) http.Handler) *Server {
	s := &Server{
		Router:      httprouter.New(),
		Middlewares: middlewares,
	}

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	NewHandler(s.Router, s.Middlewares...).ServeHTTP(w, req)
}

func (s *Server) ListenAndServe(addr string) error {
	// ListenAndServe will block
	//
	// Install our handler at the root of the standard net/http default mux.
	// This allows packages like expvar to continue working as expected.
	err := http.ListenAndServe(addr, s)
	return err
}

func (s *Server) Handle(method, p string, h http.Handler) {
	s.Router.Handler(method, p, h)
}

// POST is a shortcut for Server.Handle("POST", p, handle)
func (s *Server) Post(p string, handler http.Handler) {
	s.Handle("POST", p, handler)
}

// GET is a shortcut for Server.Handle("GET", p, handle)
func (s *Server) Get(p string, handler http.Handler) {
	s.Handle("GET", p, handler)
}

// OPTIONS is a shortcut for Server.Handle("OPTIONS", p, handle)
func (s *Server) Options(p string, handler http.Handler) {
	s.Handle("OPTIONS", p, handler)
}

// HEAD is a shortcut for Server.Handle("HEAD", p, handle)
func (s *Server) Head(p string, handler http.Handler) {
	s.Handle("HEAD", p, handler)
}

// NewHandler creates a new http handler with optional middlewares
func NewHandler(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	if len(middlewares) == 0 {
		return handler
	}

	var h http.Handler
	// give our main handler to the first middleware
	h = middlewares[len(middlewares)-1](handler)

	for i := len(middlewares) - 2; i >= 0; i-- {
		h = middlewares[i](h)
	}

	return h
}
