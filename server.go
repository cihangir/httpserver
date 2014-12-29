package httpserver

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Server is a high performance simple configurable Go HTTP server
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

// ServeHTTP implements Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	NewHandler(s.Router, s.Middlewares...).ServeHTTP(w, req)
}

// ListenAndServe serves the server with default http mux
func (s *Server) ListenAndServe(addr string) error {
	// ListenAndServe will block
	//
	// Install our handler at the root of the standard net/http default mux.
	// This allows packages like expvar to continue working as expected.
	err := http.ListenAndServe(addr, s)
	return err
}

// Handle passes handlers to the router
func (s *Server) Handle(method, p string, h http.Handler) {
	s.Router.Handler(method, p, h)
}

// Post is a shortcut for Server.Handle("POST", p, handle)
func (s *Server) Post(p string, handler http.Handler) {
	s.Handle("POST", p, handler)
}

// Get is a shortcut for Server.Handle("GET", p, handle)
func (s *Server) Get(p string, handler http.Handler) {
	s.Handle("GET", p, handler)
}

// Options is a shortcut for Server.Handle("OPTIONS", p, handle)
func (s *Server) Options(p string, handler http.Handler) {
	s.Handle("OPTIONS", p, handler)
}

// Head is a shortcut for Server.Handle("HEAD", p, handle)
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
