package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type simpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}

type server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, r *http.Request)
}

type loadBalancer struct {
	port            string
	roundRobinCount int
	servers         []server
}

func newSimpleServer(addr string) *simpleServer {
	serverUrl, err := url.Parse(addr)
	handleError(err)

	return &simpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}

}

func newLoadBalancer(port string, servers []server) *loadBalancer {
	return &loadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

func (lb *loadBalancer) getNextAvailableServer() server {
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !server.IsAlive() {
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]

	}
	lb.roundRobinCount++
	return server
}

func (lb *loadBalancer) serveProxy(rw http.ResponseWriter, r *http.Request) {
	targetServer := lb.getNextAvailableServer()
	fmt.Printf("Forwarding request to %q\n", targetServer.Address())
	targetServer.Serve(rw, r)
}

func (s *simpleServer) Address() string {
	return s.addr
}

func (s *simpleServer) IsAlive() bool {
	return true
}

func (s *simpleServer) Serve(rw http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(rw, r)
}

func main() {
	servers := []server{
		newSimpleServer("https://www.bing.com"),
		newSimpleServer("http://www.reddit.com"),
		newSimpleServer("http://www.amazon.com"),
	}

	fmt.Println(servers)

	lb := newLoadBalancer("8080", servers)

	handleRedirect := func(rw http.ResponseWriter, r *http.Request) {
		lb.serveProxy(rw, r)
	}

	http.HandleFunc("/", handleRedirect)
	fmt.Printf("Server is serving request at 'localhost:%s'\n", lb.port)
	http.ListenAndServe(":"+lb.port, nil)
}
