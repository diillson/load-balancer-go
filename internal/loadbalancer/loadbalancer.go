package loadbalancer

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type Server struct {
	URL         *url.URL
	ActiveConns int32
	Healthy     bool
}

type LoadBalancer interface {
	GetBackend() (*Server, error)
	ReleaseBackend(server *Server)
	AddBackend(server *Server)
	RemoveBackend(serverURL *url.URL)
	GetAllServers() []*Server
}

type SimpleLoadBalancer struct {
	Servers []*Server
	mux     sync.Mutex
}

const HealthCheckInterval = 10 * time.Second

func (lb *SimpleLoadBalancer) GetBackend() (*Server, error) {
	lb.mux.Lock()
	defer lb.mux.Unlock()

	if len(lb.Servers) == 0 {
		return nil, fmt.Errorf("no servers available")
	}

	// Implementação de menor conexão
	var leastConnServer *Server
	for _, server := range lb.Servers {
		// We want to make sure the server is healthy before considering it
		if server.Healthy {
			if leastConnServer == nil || atomic.LoadInt32(&server.ActiveConns) < atomic.LoadInt32(&leastConnServer.ActiveConns) {
				leastConnServer = server
			}
		}
	}

	if leastConnServer == nil {
		return nil, fmt.Errorf("no healthy servers available")
	}

	atomic.AddInt32(&leastConnServer.ActiveConns, 1)
	return leastConnServer, nil
}

func (lb *SimpleLoadBalancer) ReleaseBackend(server *Server) {
	atomic.AddInt32(&server.ActiveConns, -1)
}

func (lb *SimpleLoadBalancer) AddBackend(server *Server) {
	server.Healthy = true

	lb.mux.Lock()
	lb.Servers = append(lb.Servers, server)
	lb.mux.Unlock()
}

func (lb *SimpleLoadBalancer) RemoveBackend(serverURL *url.URL) {
	lb.mux.Lock()
	defer lb.mux.Unlock()

	for i, server := range lb.Servers {
		if server.URL.String() == serverURL.String() {
			// Remove o servidor da lista
			lb.Servers = append(lb.Servers[:i], lb.Servers[i+1:]...)
			return
		}
	}
}

func (lb *SimpleLoadBalancer) StartHealthChecks() {
	ticker := time.NewTicker(HealthCheckInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				lb.CheckAllBackendsHealth()
			}
		}
	}()
}

func (lb *SimpleLoadBalancer) CheckAllBackendsHealth() {
	for _, backend := range lb.Servers {
		backend.CheckHealth()
	}
}

func (s *Server) CheckHealth() {
	// Defina um timeout para a requisição para garantir que ela não fique pendente por muito tempo.
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Fazer uma requisição GET para o endpoint /health do servidor.
	resp, err := client.Get(s.URL.String() + "/health")

	// Se ocorrer um erro ou o status da resposta não for 200, marque o servidor como não saudável.
	if err != nil || resp.StatusCode != http.StatusOK {
		s.Healthy = false
		return
	}

	// Caso contrário, marque o servidor como saudável.
	s.Healthy = true
}

// Em loadbalancer.go
func (lb *SimpleLoadBalancer) GetAllServers() []*Server {
	lb.mux.Lock()
	defer lb.mux.Unlock()
	return lb.Servers
}
