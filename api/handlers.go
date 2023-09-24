package api

import (
	"github.com/diillson/load-balancer-go/internal/loadbalancer"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
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

	proxy := httputil.NewSingleHostReverseProxy(server.URL)

	// Adiciona/modifica o User-Agent ou qualquer header na requisição antes de ser encaminhada
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = server.URL.Scheme
		req.URL.Host = server.URL.Host
		req.URL.Path = singleJoiningSlash(server.URL.Path, req.URL.Path)

		req.Host = server.URL.Host // Defina o cabeçalho 'Host' para o do servidor alvo evitar validação de referência cruzada

		// Adicione o cabeçalho que deseja manipular ao fazer a requisição no servidor alvo
		//req.Header.Set("User-Agent", "PostmanRuntime/7.32.3")
		//req.Header.Set("Connection", "close")

		//for name, values := range req.Header {
		//	for _, value := range values {
		//		logrus.Infof("Header: %s: %s", name, value)
		//	}
		//}
	}

	proxy.ServeHTTP(c.Writer, c.Request)
	h.LB.ReleaseBackend(server)
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
