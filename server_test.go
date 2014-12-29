package httpserver

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"testing"
)

func middlewareCreator(i int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, req *http.Request) {
			p := strconv.Itoa(i)
			io.WriteString(rw, p)
			next.ServeHTTP(rw, req)
			io.WriteString(rw, p)
		}

		return http.HandlerFunc(fn)
	}
}

func startListening(s *Server) net.Listener {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}

	go http.Serve(l, s)
	return l
}

func createAddr(l net.Listener, path string) string {
	return "http://" + l.Addr().String() + path
}

func TestGlobal(t *testing.T) {
	h := New(middlewareCreator(1))
	h.Get("/1", http.HandlerFunc(helloer))

	l := startListening(h)
	defer l.Close()

	const expected = "1Hello, World!1"
	process(t, createAddr(l, "/1"), expected)

}

func TestHandler(t *testing.T) {
	h := New(middlewareCreator(1))
	h.Get("/1", NewHandler(http.HandlerFunc(helloer), middlewareCreator(2), middlewareCreator(3)))

	l := startListening(h)
	defer l.Close()

	const expected = "123Hello, World!321"
	process(t, createAddr(l, "/1"), expected)

}

func TestClean(t *testing.T) {
	h := New()
	h.Get("/1", http.HandlerFunc(helloer))

	l := startListening(h)
	defer l.Close()

	const expected = "Hello, World!"
	process(t, createAddr(l, "/1"), expected)
}

func process(t *testing.T, url string, expected string) {
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err.Error())
	}

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err.Error())
	}

	if string(r) != expected {
		t.Fatalf("expected %s, bu got %s", expected, string(r))
	}
}

func helloer(rw http.ResponseWriter, req *http.Request) {
	io.WriteString(rw, "Hello, World!")
}
