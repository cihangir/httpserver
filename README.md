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

## Install and Usage

Install the package with:

```bash
go get github.com/cihangir/httpserver
```

Import it with:

```go
import "github.com/cihangir/httpserver"
```
