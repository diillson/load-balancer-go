package api

import (
	"github.com/diillson/load-balancer-go/internal/loadbalancer"
	"github.com/diillson/load-balancer-go/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type Handler struct {
	LB loadbalancer.LoadBalancer
}

type ServerInput struct {
	URLString   string `json:"URL"`
	ActiveConns int32  `json:"ActiveConns"`
	Healthy     bool   `json:"Healthy"`
}

func (h *Handler) ProxyHandler(c *gin.Context) {
	server, err := h.LB.GetBackend()
	if err != nil {
		logrus.Error("Failed to get backend:", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service unavailable"})
		return
	}

	defer h.LB.ReleaseBackend(server) // Garantir que a conexão seja liberada ao finalizar a requisição

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = server.URL.Scheme
			req.URL.Host = server.URL.Host
			req.URL.Path = singleJoiningSlash(server.URL.Path, req.URL.Path)

			// Defina o cabeçalho 'Host' para o do servidor alvo para evitar validação de referência cruzada
			req.Host = server.URL.Host

			// Adicione o cabeçalho que deseja manipular ao fazer a requisição no servidor alvo
			//req.Header.Set("User-Agent", "PostmanRuntime/7.32.3")
			//req.Header.Set("Connection", "close")
		},
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		ErrorLog: logging.GetLogrusAdapter(), // para integrar com o logger do logrus
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

// função `singleJoiningSlash` para corrigir os paths
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func (h *Handler) AddServerHandler(c *gin.Context) {
	var input ServerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	serverURL, err := url.Parse(input.URLString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return
	}

	server := &loadbalancer.Server{
		URL:         serverURL,
		ActiveConns: input.ActiveConns,
		Healthy:     input.Healthy,
	}
	h.LB.AddBackend(server)
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *Handler) RemoveServerHandler(c *gin.Context) {
	var input ServerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	serverURL, err := url.Parse(input.URLString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return
	}

	h.LB.RemoveBackend(serverURL)
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *Handler) HealthCheckHandler(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func (h *Handler) ListServersHandler(c *gin.Context) {
	servers := h.LB.GetAllServers()

	c.JSON(http.StatusOK, servers)
}
