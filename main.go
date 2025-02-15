package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type LoadBalancer struct {
	pool 			*ServerPool
	algorithm 		func(servers []*Server) *Server
	interval 		time.Duration
}

type Server struct {
	Addr 		string
	Weight 		int
	Connections int64
}

type ServerPool struct {
	servers sync.Map
}

func (sp *ServerPool) AddServer(addr string, weight int) {
	sp.servers.Store(addr, &Server{Addr: addr, Weight: weight})
}

func (sp *ServerPool) RemoveServer(addr string) {
	sp.servers.Delete(addr)
}

func (sp *ServerPool) GetServers() []*Server {
	servers := make([]*Server, 0)
	sp.servers.Range(func(key, value interface{}) bool {
		servers = append(servers, value.(*Server))
		return true
	})

	return servers
}


func LeastConnections(servers []*Server) *Server {
	var least *Server
	for _, server := range servers {
		if least == nil || server.Connections < least.Connections {
			least = server
		}
	}

	return least
}


func healthCheck(server *Server) bool {
	url, err := url.Parse("http://" + server.Addr)
	
	if err != nil { //server down
		fmt.Println("Server down", server.Addr)
		return false
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url.String())

	if err != nil || resp.StatusCode != 200 {
		fmt.Println("Server down", server.Addr)
		return false
	}
	defer resp.Body.Close()

	return true
}


func performHealthCheck(sp *ServerPool, interval time.Duration) {
	for {
		time.Sleep(interval)
		sp.servers.Range(func(key, value interface{}) bool {
			server := value.(*Server)
			if !healthCheck(server) {
				fmt.Println("Removing server", server.Addr)
				sp.RemoveServer(server.Addr) //this is bad, if this server is healthy again, it will not be added back
			}
			return true
		})
	}
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	servers := lb.pool.GetServers()
	if len(servers) == 0 {
		http.Error(w, "No server available", http.StatusServiceUnavailable)

	}
	server := lb.algorithm(lb.pool.GetServers())

	atomic.AddInt64(&server.Connections, 1)
	defer func() {
		atomic.AddInt64(&server.Connections, -1)
	}()

	httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: server.Addr}).ServeHTTP(w, r)
}


func main() {
	sp := &ServerPool{}
	sp.AddServer("127.0.0.1:8081", 1)
	sp.AddServer("127.0.0.1:8082", 2)
	sp.AddServer("127.0.0.1:8083", 3)

	lb := &LoadBalancer{
		pool: sp,
		algorithm: LeastConnections,
		interval: 5 * time.Second,
	}

	go performHealthCheck(sp, 5 * time.Second) //run healthcheck on a seperate thread.

	http.ListenAndServe(":8080", lb)
	fmt.Println("Load balancer started at :8080")
}
