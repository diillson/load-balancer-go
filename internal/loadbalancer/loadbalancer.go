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
	server.Healthy = false

	lb.mux.Lock()
	lb.Servers = append(lb.Servers, server)
	lb.mux.Unlock()
	go server.CheckHealth()
}

func (lb *SimpleLoadBalancer) RemoveBackend(serverURL *url.URL) {
	lb.mux.Lock()
	defer lb.mux.Unlock()

	for i, server := range lb.Servers {
		if server.URL.String() == serverURL.String() {
			atomic.StoreInt32(&server.ActiveConns, 0) // Resetting active connections for removed server
			// Remove o servidor da lista
			lb.Servers = append(lb.Servers[:i], lb.Servers[i+1:]...)
			return
		}
	}
}

func (lb *SimpleLoadBalancer) StartHealthChecks() {
	ticker := time.NewTicker(HealthCheckInterval)
	go func() {
		for range ticker.C {
			lb.CheckAllBackendsHealth()
		}
	}()
}

func (lb *SimpleLoadBalancer) CheckAllBackendsHealth() {
	for _, backend := range lb.Servers {
		go backend.CheckHealth()
	}
}

func (lb *SimpleLoadBalancer) Initialize() {
	// Iniciar verificações de saúde.
	lb.StartHealthChecks()

	// Verificar a saúde inicial de todos os servidores.
	lb.CheckAllBackendsHealth()
}

func (s *Server) CheckHealth() {
	// Defina um timeout para a requisição para garantir que ela não fique pendente por muito tempo.
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:      10,
			IdleConnTimeout:   30 * time.Second,
			DisableKeepAlives: true,
			MaxConnsPerHost:   10,
		},
	}

	// Fazer uma requisição GET para o endpoint /health do servidor.
	resp, err := client.Get(s.URL.String() + "/health")

	// Se ocorrer um erro ou o status da resposta não for 200, marque o servidor como não saudável.
	if err != nil || (resp.StatusCode != http.StatusOK && resp.StatusCode < http.StatusInternalServerError) {
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
