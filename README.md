[![GoDoc](https://godoc.org/github.com/cihangir/httpserver?status.svg)](https://godoc.org/github.com/cihangir/httpserver)
[![Build Status](https://travis-ci.org/cihangir/httpserver.svg)](https://travis-ci.org/cihangir/httpserver)

httpserver
==========

A high performance simple configurable Go HTTP server that is compatible with http.Handler

```go
package main

import "github.com/cihangir/httpserver"

func middleware(next http.Handler) http.Handler {
    fn := func(rw http.ResponseWriter, req *http.Request) {
        // do first
        next.ServeHTTP(rw, req)
        // do then
    }

    return http.HandlerFunc(fn)
}

func main() {
    s := httpserver.New(middleware)
    s.Get("/1", helloer)
    s.Get("/2", httpserver.NewHandler(
        http.HandlerFunc(helloer),
        middleware2,
        middleware3,
    ))
    s.ListenAndServe(addr)
}

func helloer(rw http.ResponseWriter, req *http.Request) {
    io.WriteString(rw, "Hello, World!")
}

```

##Middlewares

Middleware are just http.HandlerFunc's

```
func middlewareCreator(next http.Handler) http.Handler {
    fn := func(rw http.ResponseWriter, req *http.Request) {
        p := "my middleware"

        io.WriteString(rw, p)
        next.ServeHTTP(rw, req)
        io.WriteString(rw, p)
    }

    return http.HandlerFunc(fn)
}

```


Execution of the handlers are not bi-directional, you can wrap all handler with the first middleware. Following is the order of execution;
```
//0 START
//1 START
//2 START
//2 END
//1 END
//0 END
```

## Features
**Simple:** Compatible with stdlib's handler structure

**Pluggable:** You can have global and handler based middlewares

**Performant:** It uses [julienschmidt/httprouter](https://github.com/julienschmidt/go-http-routing-benchmark) as router

## Install

Install the package with:

```bash
go get github.com/cihangir/httpserver
```

Import it with:

```go
import "github.com/cihangir/httpserver"
```

## Usage

#### func  NewHandler

```go
func NewHandler(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler
```
NewHandler creates a new http handler with optional middlewares

#### type Server

```go
type Server struct {
	// Router holds a performant router for requests
	Router *httprouter.Router

	// Middlewares holds global middlewares for the server
	Middlewares []func(next http.Handler) http.Handler
}
```

Server is a high performance simple configurable Go HTTP server

#### func  New

```go
func New(middlewares ...func(http.Handler) http.Handler) *Server
```
New creates a server

#### func (*Server) Get

```go
func (s *Server) Get(p string, handler http.Handler)
```
Get is a shortcut for Server.Handle("GET", p, handle)

#### func (*Server) Handle

```go
func (s *Server) Handle(method, p string, h http.Handler)
```
Handle passes handlers to the router

#### func (*Server) Head

```go
func (s *Server) Head(p string, handler http.Handler)
```
Head is a shortcut for Server.Handle("HEAD", p, handle)

#### func (*Server) ListenAndServe

```go
func (s *Server) ListenAndServe(addr string) error
```
ListenAndServe serves the server with default http mux

#### func (*Server) Options

```go
func (s *Server) Options(p string, handler http.Handler)
```
Options is a shortcut for Server.Handle("OPTIONS", p, handle)

#### func (*Server) Post

```go
func (s *Server) Post(p string, handler http.Handler)
```
Post is a shortcut for Server.Handle("POST", p, handle)

#### func (*Server) ServeHTTP

```go
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request)
```
ServeHTTP implements Handler interface
